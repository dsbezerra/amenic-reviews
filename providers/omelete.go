package providers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dsbezerra/amenic-reviews/extractors"
	"github.com/dsbezerra/amenic-reviews/helpers"
	"github.com/dsbezerra/amenic-reviews/stringutil"
)

const (
	omeleteBaseURL = "https://omelete.com.br"

	omeleteSectionTypeMovies = "section_type:movies"
	omeleteSectionTypeShows  = "section_type:shows"
)

type (
	Omelete struct {
		BaseURL   string
		Selectors OmeleteSelectors
	}

	OmeleteProperties struct {
		Author     string
		Title      string
		Subtitle   string
		ReviewBody string
	}

	OmeleteSelectors struct {
		SearchTotalCount string
		Review           OmeleteReviewSelectors
		ReviewsSections  OmeleteReviewsSectionSelectors
	}

	OmeleteReviewsSectionSelectors struct {
		Container string
		Label     string
		ListItem  string
		ImageURL  string
		Title     string
		Subtitle  string
		Date      string
	}

	OmeleteReviewSelectors struct {
		Container         string
		RatingDescription string
	}

	// OmeleteReviewsSection represents a review section
	OmeleteReviewsSection struct {
		Omelete     *Omelete                    `json:"-"`
		Label       string                      `json:"label"`
		SectionType string                      `json:"section_type"`
		Results     []OmeleteReviewSearchResult `json:"results"`
		Pagination  Pagination                  `json:"pagination"`
	}

	// OmeleteReviewSearchResult represents a review search result
	OmeleteReviewSearchResult struct {
		Omelete  *Omelete               `json:"-"`
		Section  *OmeleteReviewsSection `json:"-"`
		ImageURL string                 `json:"image_url"`
		Date     *time.Time             `json:"publish_date"`
		Type     string                 `json:"type"`
		Title    string                 `json:"title"`
		Subtitle string                 `json:"subtitle"`
		Path     string                 `json:"path"`
	}

	// OmeleteReview represents a review page
	OmeleteReview struct {
		Headline          string
		Description       string
		Author            string
		PublishDate       string
		Content           string
		ContentHtml       string
		Rating            int
		RatingDescription string
		ReviewURL         string
	}
)

// NewOmeleteProvider creates a new omelete provider
func NewOmeleteProvider() *Omelete {
	result := &Omelete{
		BaseURL: omeleteBaseURL,
		Selectors: OmeleteSelectors{
			SearchTotalCount: "#filters > div.principal > ul > li.item.active > span > a > span.count",
			ReviewsSections: OmeleteReviewsSectionSelectors{
				Container: "#conteudo",
				Label:     "h1",
				ListItem:  ".include.search-content-type",
				ImageURL:  "img",
				Title:     ".title",
				Subtitle:  ".subtitle",
				Date:      "span.date",
			},
			Review: OmeleteReviewSelectors{
				Container:         "div.article-main",
				RatingDescription: "div.rating-ficha > span.nota-texto",
			},
		},
	}

	return result
}

// GetMovieReviewInfo retorna crítica de um determinado filme
func (o *Omelete) GetMovieReviewInfo(path string) (*OmeleteReview, error) {
	return o.getReviewInfo(path, omeleteSectionTypeMovies)
}

// GetShowReviewInfo retorna crítica de um determinado filme
func (o *Omelete) GetShowReviewInfo(path string) (*OmeleteReview, error) {
	return o.getReviewInfo(path, omeleteSectionTypeShows)
}

// SearchMovieReviews encontra as críticas de filmes que satisfazem a query especificada
func (o *Omelete) SearchMovieReviews(query string) ([]OmeleteReviewSearchResult, error) {
	return o.searchReviewsFor(query, omeleteSectionTypeMovies, 1)
}

// SearchShowsReviews encontra as críticas de séries que satisfazem a query especificada
func (o *Omelete) SearchShowsReviews(query string) ([]OmeleteReviewSearchResult, error) {
	return o.searchReviewsFor(query, omeleteSectionTypeShows, 1)
}

// GetMovieReviewsSection Retorna a primeira page da seção de críticas de filmes
func (o *Omelete) GetMovieReviewsSection() (*OmeleteReviewsSection, error) {
	return o.getReviewsSection(omeleteSectionTypeMovies, 1)
}

// GetShowsReviewsSection Retorna a primeira page da seção de críticas de séries de TV
func (o *Omelete) GetShowsReviewsSection() (*OmeleteReviewsSection, error) {
	return o.getReviewsSection(omeleteSectionTypeShows, 1)
}

// PreviousPage Retorna a página de críticas anterior da seção
func (s *OmeleteReviewsSection) PreviousPage() (*OmeleteReviewsSection, error) {
	if s.Pagination.PreviousPage > 1 {
		return s.Omelete.getReviewsSection(s.SectionType, s.Pagination.PreviousPage)
	}

	return nil, errors.New("page not found")
}

// NextPage Retorna a próxima página de críticas da seção
func (s *OmeleteReviewsSection) NextPage() (*OmeleteReviewsSection, error) {
	return s.Omelete.getReviewsSection(s.SectionType, s.Pagination.NextPage)
}

// GetReviewInfo Returns a pointer to a OmeleteReview that matches the given path
func (o *Omelete) getReviewInfo(path string, sectionType string) (*OmeleteReview, error) {
	if path == "" {
		return nil, errors.New("path is missing")
	}

	finalURL := o.BaseURL + getReviewInfoPathFor(sectionType)
	if !strings.HasPrefix(path, "/") {
		finalURL += "/"
	}
	finalURL += path
	fmt.Printf("URL: %s\n", finalURL)
	doc, err := helpers.NewDocument(finalURL, "utf-8")
	if err != nil {
		return nil, err
	}

	result := &OmeleteReview{}
	result.ReviewURL = finalURL
	o.handleReviewInfoResponse(doc, result)
	return result, nil
}

// GetReviewInfo Retorna a crítica do resultado (seja filme ou série ou qualquer outro implementado)
func (r *OmeleteReviewSearchResult) GetReviewInfo() (*OmeleteReview, error) {
	finalURL := r.Omelete.BaseURL + r.Path
	doc, err := helpers.NewDocument(finalURL, "utf-8")
	if err != nil {
		return nil, err
	}

	result := &OmeleteReview{}
	result.ReviewURL = finalURL
	r.Omelete.handleReviewInfoResponse(doc, result)
	return result, nil
}

func getReviewsSearchPathFor(query string, sectionType string, page int) string {
	if query == "" {
		return ""
	}

	result := "/busca/?q=" + strings.Replace(query, " ", "+", -1)
	section := "&secao="
	switch sectionType {
	case omeleteSectionTypeMovies:
		section += "filmes"
	case omeleteSectionTypeShows:
		section += "series-tv"
	default:
	}

	result += section + "&tipo=critica&pagina=" + strconv.Itoa(page)
	return result
}

func (o *Omelete) searchReviewsFor(query string, sectionType string, page int) ([]OmeleteReviewSearchResult, error) {
	finalURL := o.BaseURL + getReviewsSearchPathFor(query, sectionType, page)

	doc, err := helpers.NewDocument(finalURL, "utf-8")
	if err != nil {
		return nil, err
	}
	// NOTE(diego): We use handleReviewsSectionResponse here since this request will return
	// almost the same page handled in reviews section
	section := &OmeleteReviewsSection{
		Omelete:     o,
		SectionType: sectionType,
	}
	o.handleReviewsSectionResponse(doc, section)
	// TODO(diego): Visit all pages if necessary
	return section.Results, nil
}

func getSectionPathFor(sectionType string) string {
	result := ""

	switch sectionType {
	case omeleteSectionTypeMovies:
		result = "/filmes" + result
	case omeleteSectionTypeShows:
		result = "/series-tv" + result
	default:
		result = "/home" + result
	}

	return result
}

func getReviewsSectionPathFor(sectionType string, page int) string {
	result := "/critica"
	result += getSectionPathFor(sectionType)
	if page > 1 {
		result += "/?pagina=" + strconv.Itoa(page)
	}
	return result
}

func getReviewInfoPathFor(sectionType string) string {
	result := getSectionPathFor(sectionType)
	result += "/criticas"
	return result
}

func (o *Omelete) getReviewsSection(sectionType string, page int) (*OmeleteReviewsSection, error) {
	finalURL := o.BaseURL + getReviewsSectionPathFor(sectionType, page)
	result := &OmeleteReviewsSection{
		Omelete:     o,
		SectionType: sectionType,
	}

	doc, err := helpers.NewDocument(finalURL, "utf-8")
	if err != nil {
		return nil, err
	}
	o.handleReviewsSectionResponse(doc, result)
	return result, nil
}

func (o *Omelete) handleReviewsSectionResponse(doc *goquery.Document, section *OmeleteReviewsSection) {
	selectors := o.Selectors.ReviewsSections
	container := doc.Find(selectors.Container)
	if container != nil && container.Length() != 0 {
		section.Label = extractors.GetText(container, selectors.Label).Trimmed
		items := container.Find(selectors.ListItem)
		items.Each(func(i int, s *goquery.Selection) {
			item := OmeleteReviewSearchResult{
				Omelete:  o,
				Section:  section,
				Type:     section.SectionType,
				Title:    extractors.GetText(s, selectors.Title).Trimmed,
				Subtitle: extractors.GetText(s, selectors.Subtitle).Trimmed,
			}

			// Retrieve image url.
			image := extractors.GetImage(s, selectors.ImageURL).SrcURL

			// TODO(diego): Helper to get date.
			date := extractors.GetText(s, selectors.Date)
			if date != nil {
				parts := date.Split("|")
				if len(parts) == 2 {
					// Get date.
					ate := stringutil.EatUntilAlpha(parts[0].Trimmed)
					t, err := createDateFromString(ate, "/", false)
					if err != nil {
						// TODO: Proper logging
						fmt.Println(err.Error())
					} else {
						// Apply time
						h, m := stringutil.BreakByToken(strings.TrimSpace(parts[1].Trimmed), 'h')
						if h != "" && m != "" {
							hours, _ := strconv.Atoi(h)
							minutes, _ := strconv.Atoi(m)
							t = t.Add(time.Hour * time.Duration(hours))
							t = t.Add(time.Minute * time.Duration(minutes))
						}

						item.Date = &t
					}
				}
			}

			// Retrieve review url path.
			href, exists := s.Find("a").Attr("href")
			if !exists {
				fmt.Printf("Couldn't find href for item: %v\n", item)
			} else {
				// Remove key param from path
				lastIndex := strings.LastIndex(href, "?key=")
				if lastIndex > -1 {
					href = href[0:lastIndex]
				}
			}
			item.ImageURL = image
			item.Path = href

			if item.Title != "" && item.Subtitle != "" && item.Path != "" {
				section.Results = append(section.Results, item)
			}
		})
	}

	currentPage := container.Find(".centered-paginator > ul > li.active")
	value, exists := currentPage.Attr("data-value")
	if exists {
		valueInt, err := strconv.Atoi(value)
		if err == nil {
			section.Pagination.CurrentPage = valueInt
			section.Pagination.NextPage = valueInt + 1
			if valueInt > 1 {
				section.Pagination.PreviousPage = valueInt - 1
			}
		}
	}
}

func (o *Omelete) handleReviewInfoResponse(doc *goquery.Document, review *OmeleteReview) {
	selectors := o.Selectors.Review
	container := doc.Find(selectors.Container)
	if container != nil && container.Length() != 0 {
		if review == nil {
			review = &OmeleteReview{}
		}

		review.Headline = extractors.GetHeadline(container).Trimmed
		review.Description = extractors.GetDescription(container).Trimmed
		review.PublishDate = extractors.GetItemProperty(container, "datePublished").ContentText.Trimmed
		review.Author = extractors.GetAuthor(container).Trimmed

		reviewBody := extractors.GetItemProperty(container, "reviewBody")
		if reviewBody != nil {
			review.ContentHtml = reviewBody.Content
			review.Content = reviewBody.ContentText.Trimmed
		}

		review.Rating = extractors.GetItemProperty(container, "ratingValue").ContentText.AsInt()
		desc := extractors.GetText(container, selectors.RatingDescription).Trimmed
		if desc != "" {
			review.RatingDescription = stringutil.SubstringBetween(desc, '(', ')')
		}
	}
}

func createDateFromString(s string, delim string, hasYear bool) (time.Time, error) {
	var result time.Time
	if s == "" || delim == "" {
		return result, errors.New("Date string or delimiter must be specified")
	}

	parts := strings.Split(s, delim)
	day, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	month, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

	var year int
	if hasYear {
		year, _ = strconv.Atoi(parts[2])

		if len(parts[2]) == 2 {
			// @DumbHack
			year += 2000
		}

	} else {
		year = time.Now().Year()
	}

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.UTC
	}
	result = time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	return result, nil
}

func getOmeleteReviewSelectors() OmeleteReviewSelectors {
	return OmeleteReviewSelectors{
		Container:         "div.article-main",
		RatingDescription: "div.rating-ficha > span.nota-texto",
	}
}
