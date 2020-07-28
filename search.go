package reddit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SearchService handles communication with the search
// related methods of the Reddit API.
//
// For searches to include NSFW results, the user must
// enable the following setting in their preferences:
// "include not safe for work (NSFW) search results in searches"
// Note: The "limit" parameter in searches is prone to inconsistent
// behaviour, e.g. sometimes limit=1 returns nothing when it should.
//
// Reddit API docs: https://www.reddit.com/dev/api/#section_search
type SearchService struct {
	client *Client
}

func newSearchOptions(opts ...SearchOptionSetter) url.Values {
	searchOptions := make(url.Values)
	for _, opt := range opts {
		opt(searchOptions)
	}
	return searchOptions
}

// SearchOptionSetter sets values for the options.
type SearchOptionSetter func(opts url.Values)

// SetAfter sets the after option.
func SetAfter(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("after", v)
	}
}

// SetBefore sets the before option.
func SetBefore(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("before", v)
	}
}

// SetLimit sets the limit option.
// Warning: It seems like setting the limit to 1 sometimes returns 0 results.
func SetLimit(v int) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("limit", fmt.Sprint(v))
	}
}

// SortByHot sets the sort option to return the hottest results first.
func SortByHot(opts url.Values) {
	opts.Set("sort", "hot")
}

// SortByBest sets the sort option to return the best results first.
func SortByBest(opts url.Values) {
	opts.Set("sort", "best")
}

// SortByNew sets the sort option to return the newest results first.
func SortByNew(opts url.Values) {
	opts.Set("sort", "new")
}

// SortByRising sets the sort option to return the rising results first.
func SortByRising(opts url.Values) {
	opts.Set("sort", "rising")
}

// SortByControversial sets the sort option to return the most controversial results first.
func SortByControversial(opts url.Values) {
	opts.Set("sort", "controversial")
}

// SortByTop sets the sort option to return the top results first.
func SortByTop(opts url.Values) {
	opts.Set("sort", "top")
}

// SortByRelevance sets the sort option to return the most relevant results first.
func SortByRelevance(opts url.Values) {
	opts.Set("sort", "relevance")
}

// SortByNumberOfComments sets the sort option to return the results with the highest
// number of comments first.
func SortByNumberOfComments(opts url.Values) {
	opts.Set("sort", "comments")
}

// FromThePastHour sets the timespan option to return results from the past hour.
func FromThePastHour(opts url.Values) {
	opts.Set("t", "hour")
}

// FromThePastDay sets the timespan option to return results from the past day.
func FromThePastDay(opts url.Values) {
	opts.Set("t", "day")
}

// FromThePastWeek sets the timespan option to return results from the past week.
func FromThePastWeek(opts url.Values) {
	opts.Set("t", "week")
}

// FromThePastMonth sets the timespan option to return results from the past month.
func FromThePastMonth(opts url.Values) {
	opts.Set("t", "month")
}

// FromThePastYear sets the timespan option to return results from the past year.
func FromThePastYear(opts url.Values) {
	opts.Set("t", "year")
}

// FromAllTime sets the timespan option to return results from all time.
func FromAllTime(opts url.Values) {
	opts.Set("t", "all")
}

// setType sets the type option.
// It could be user, link, sr (subreddit).
// For mod actions, it's for the type of action (e.g. "banuser", "spamcomment").
func setType(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("type", v)
	}
}

// setQuery sets the q option.
func setQuery(v string) SearchOptionSetter {
	return func(opts url.Values) {
		opts.Set("q", v)
	}
}

// setRestrict sets the restrict_sr option.
func setRestrict(opts url.Values) {
	opts.Set("restrict_sr", "true")
}

// Posts searches for posts.
// If the list of subreddits provided is empty, the search is run against r/all.
func (s *SearchService) Posts(ctx context.Context, query string, subreddits []string, opts ...SearchOptionSetter) (*Posts, *Response, error) {
	opts = append(opts, setType("link"), setQuery(query))

	path := "search"
	if len(subreddits) > 0 {
		path = fmt.Sprintf("r/%s/search", strings.Join(subreddits, "+"))
		opts = append(opts, setRestrict)
	}

	form := newSearchOptions(opts...)
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getPosts(), resp, nil
}

// Subreddits searches for subreddits.
// The sort and timespan options don't affect the results for this search.
func (s *SearchService) Subreddits(ctx context.Context, query string, opts ...SearchOptionSetter) (*Subreddits, *Response, error) {
	opts = append(opts, setType("sr"), setQuery(query))
	form := newSearchOptions(opts...)

	path := "search"
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getSubreddits(), resp, nil
}

// Users searches for users.
// The sort and timespan options don't affect the results for this search.
func (s *SearchService) Users(ctx context.Context, query string, opts ...SearchOptionSetter) (*Users, *Response, error) {
	opts = append(opts, setType("user"), setQuery(query))
	form := newSearchOptions(opts...)

	path := "search"
	path = addQuery(path, form)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(rootListing)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.getUsers(), resp, nil
}
