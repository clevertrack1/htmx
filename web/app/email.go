package app

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var sessionStore sync.Map

type EmailApp struct {
	Tmpl *template.Template
}

type EmailSession struct {
	Emails []Email
}

type Email struct {
	Id          int
	Sender      string
	Description string
	Body        string
	Starred     bool
	Archived    bool
}

type EmailRender struct {
	Emails        []Email
	SelectedEmail Email
	CurrentView   string
}

func (a *EmailApp) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /emails", a.renderEmails)
	mux.HandleFunc("GET /emails/{id}", a.renderEmail)
	mux.HandleFunc("POST /emails/{id}/star", a.starEmail)
	mux.HandleFunc("POST /emails/{id}/archive", a.archiveEmail)
	mux.HandleFunc("POST /emails/search", a.searchEmails)
}

// renderEmails adapts to the requested view (inbox, archived, starred).
func (a *EmailApp) renderEmails(w http.ResponseWriter, r *http.Request) {
	s := ensureSession(w, r)

	view := r.URL.Query().Get("view") // e.g. "archived" or "starred"

	var visibleEmails []Email
	switch view {
	case "archived":
		// Show only archived emails
		visibleEmails = filterArchived(s.Emails, true)
	case "starred":
		// Show only starred (and not archived)
		visibleEmails = filterStarred(s.Emails, true)
	default:
		// Default to inbox view (all non-archived)
		view = "inbox"
		visibleEmails = filterArchived(s.Emails, false)
	}

	data := EmailRender{
		Emails:        visibleEmails,
		CurrentView:   view,
		SelectedEmail: Email{Id: -1},
	}
	if len(visibleEmails) > 0 {
		data.SelectedEmail = visibleEmails[0]
	}
	a.Tmpl.ExecuteTemplate(w, "emailview", data)
}

func (a *EmailApp) renderEmail(w http.ResponseWriter, r *http.Request) {
	s := ensureSession(w, r)

	uniqueID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	emailPtr := findEmailByID(s.Emails, uniqueID)
	if emailPtr == nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}

	// Just assume we want to keep the same `view` as on the listing page:
	view := r.URL.Query().Get("view")
	if view == "" {
		view = "inbox" // fallback
	}

	// Filter the sidebar list
	var visibleEmails []Email
	switch view {
	case "archived":
		visibleEmails = filterArchived(s.Emails, true)
	case "starred":
		visibleEmails = filterStarred(s.Emails, true)
	default:
		visibleEmails = filterArchived(s.Emails, false)
		view = "inbox"
	}

	data := EmailRender{
		Emails:        visibleEmails,
		SelectedEmail: *emailPtr,
		CurrentView:   view,
	}
	a.Tmpl.ExecuteTemplate(w, "emailview", data)
}

func (a *EmailApp) starEmail(w http.ResponseWriter, r *http.Request) {
	s := ensureSession(w, r)
	uniqueID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	emailPtr := findEmailByID(s.Emails, uniqueID)
	if emailPtr == nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}
	emailPtr.Starred = !emailPtr.Starred

	// Rerender with the same 'view' param, e.g. ?view=starred
	view := r.URL.Query().Get("view")
	var visibleEmails []Email
	switch view {
	case "archived":
		visibleEmails = filterArchived(s.Emails, true)
	case "starred":
		visibleEmails = filterStarred(s.Emails, true)
	default:
		view = "inbox"
		visibleEmails = filterArchived(s.Emails, false)
	}

	// Determine the selected email in the new filtered list
	selected := Email{Id: -1}
	if e := findEmailByID(visibleEmails, emailPtr.Id); e != nil {
		selected = *e
	} else if len(visibleEmails) > 0 {
		selected = visibleEmails[0]
	}

	data := EmailRender{
		Emails:        visibleEmails,
		SelectedEmail: selected,
		CurrentView:   view,
	}
	a.Tmpl.ExecuteTemplate(w, "emailview", data)
}

func (a *EmailApp) archiveEmail(w http.ResponseWriter, r *http.Request) {
	s := ensureSession(w, r)
	uniqueID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	emailPtr := findEmailByID(s.Emails, uniqueID)
	if emailPtr == nil {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	}
	emailPtr.Archived = !emailPtr.Archived

	view := r.URL.Query().Get("view")
	var visibleEmails []Email
	switch view {
	case "archived":
		visibleEmails = filterArchived(s.Emails, true)
	case "starred":
		visibleEmails = filterStarred(s.Emails, true)
	default:
		view = "inbox"
		visibleEmails = filterArchived(s.Emails, false)
	}

	// Determine the selected email in the new filtered list
	selected := Email{Id: -1}
	if e := findEmailByID(visibleEmails, emailPtr.Id); e != nil {
		selected = *e
	} else if len(visibleEmails) > 0 {
		selected = visibleEmails[0]
	}

	data := EmailRender{Emails: visibleEmails, SelectedEmail: selected, CurrentView: view}
	a.Tmpl.ExecuteTemplate(w, "emailview", data)
}

func (a *EmailApp) searchEmails(w http.ResponseWriter, r *http.Request) {
	s := ensureSession(w, r)

	// Grab the query from form data. If you name the field "q" in the <input>, use FormValue("q").
	searchTerm := r.FormValue("searchQuery")

	view := r.URL.Query().Get("view")
	var visibleEmails []Email
	switch view {
	case "archived":
		visibleEmails = filterArchived(s.Emails, true)
	case "starred":
		visibleEmails = filterStarred(s.Emails, true)
	default:
		view = "inbox"
		visibleEmails = filterArchived(s.Emails, false)
	}

	// Filter emails by matching subject, body, sender, etc. (Adjust the logic as you see fit.)
	results := filterEmailsBySearch(visibleEmails, searchTerm)

	// Render the "emailsList" partial, which expects .Emails and .SelectedEmail, etc.
	data := EmailRender{
		Emails:        results,
		CurrentView:   "", // or "search" if you want to treat it as a special view
		SelectedEmail: Email{Id: -1},
	}

	// Execute ONLY the sub-template that lists emails, without the entire layout
	a.Tmpl.ExecuteTemplate(w, "emailList", data)
}

// Utility helpers:

// Simple search filter, e.g. checking if searchTerm is in Subject, Body, or Sender
func filterEmailsBySearch(emails []Email, searchTerm string) []Email {
	if searchTerm == "" {
		return emails
	}
	var result []Email
	lowerQuery := strings.ToLower(searchTerm)
	for _, e := range emails {
		if strings.Contains(strings.ToLower(e.Description), lowerQuery) ||
			strings.Contains(strings.ToLower(e.Body), lowerQuery) ||
			strings.Contains(strings.ToLower(e.Sender), lowerQuery) {
			result = append(result, e)
		}
	}
	return result
}

func findEmailByID(emails []Email, id int) *Email {
	for i := range emails {
		if emails[i].Id == id {
			return &emails[i]
		}
	}
	return nil
}

// Optionally filter starred + not archived or archived + starred, as you see fit.
func filterStarred(emails []Email, wantStarred bool) []Email {
	var result []Email
	for _, e := range emails {
		if e.Starred == wantStarred && !e.Archived {
			result = append(result, e)
		}
	}
	return result
}

func filterArchived(emails []Email, wantArchived bool) []Email {
	var result []Email
	for _, e := range emails {
		if e.Archived == wantArchived {
			result = append(result, e)
		}
	}
	return result
}

func ensureSession(w http.ResponseWriter, r *http.Request) *EmailSession {
	c, err := r.Cookie("session_id")
	if err != nil || c == nil || c.Value == "" {
		sid := generateSessionID()
		s := &EmailSession{}
		s.Emails = defaultEmails

		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: sid,
			Path:  "/",
		})
		return s
	}

	sid := c.Value
	v, ok := sessionStore.Load(sid)
	if !ok {
		s := &EmailSession{}
		s.Emails = defaultEmails
		sessionStore.Store(sid, s)
		return s
	}

	return v.(*EmailSession)
}

func generateSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback-session-id"
	}
	return hex.EncodeToString(b)
}

var defaultEmails = []Email{
	{
		Id:          0,
		Sender:      "alice.adams@company.com",
		Description: "Project Kickoff",
		Body:        "Hello Team,\n\nPlease join me for our official project kickoff meeting tomorrow at 9 AM. Looking forward to collaborating with you!\n\nBest,\nAlice",
	},
	{
		Id:          1,
		Sender:      "bob.barker@company.com",
		Description: "Weekly Sprint Retro",
		Body:        "Hi Everyone,\n\nJust a reminder that our sprint retrospective is scheduled for 3 PM today. Bring any wild ideas you’d like to share!\n\nThanks,\nBob",
	},
	{
		Id:          2,
		Sender:      "carol.connors@company.com",
		Description: "Budget Approval",
		Body:        "Team,\n\nGreat news! The budget for Q1 has been approved. Let’s stay on track and utilize our resources wisely.\n\nRegards,\nCarol",
	},
	{
		Id:          3,
		Sender:      "david.daniels@company.com",
		Description: "Upcoming Hackathon",
		Body:        "Hello Folks,\n\nWe’re hosting a company-wide hackathon next Friday. Come prepared to build, break, and have fun!\n\nCheers,\nDavid",
	},
	{
		Id:          4,
		Sender:      "eve.evans@company.com",
		Description: "Employee Feedback",
		Body:        "Dear Team,\n\nWe value your feedback. If you have any suggestions for improvements, please let me know.\n\nSincerely,\nEve",
	},
	{
		Id:          5,
		Sender:      "frank.fields@company.com",
		Description: "Team Building Event",
		Body:        "Hey Everyone,\n\nI’m planning a fun team outing next month. Let me know if you have any cool activity ideas!\n\nBest,\nFrank",
	},
	{
		Id:          6,
		Sender:      "grace.garcia@company.com",
		Description: "Quarterly Check-In",
		Body:        "Hi,\n\nIt’s time for our quarterly performance check-in. Please schedule a short meeting at your earliest convenience.\n\nThank you,\nGrace",
	},
	{
		Id:          7,
		Sender:      "henry.hughes@company.com",
		Description: "Marketing Updates",
		Body:        "Team,\n\nWe have new marketing initiatives launching next week. Stay tuned for more details.\n\nRegards,\nHenry",
	},
	{
		Id:          8,
		Sender:      "irene.ingersoll@company.com",
		Description: "Client Onboarding",
		Body:        "Hello,\n\nWe have a client onboarding session next Monday. Please review the attached materials before joining.\n\nThanks,\nIrene",
	},
	{
		Id:          9,
		Sender:      "jack.jenkins@company.com",
		Description: "Team Kudos",
		Body:        "Hello Crew,\n\nWanted to give kudos for everyone’s hard work this quarter. Your dedication is truly appreciated.\n\nAll the best,\nJack",
	},
	{
		Id:          10,
		Sender:      "karen.knox@company.com",
		Description: "Security Policy Reminder",
		Body:        "Hi All,\n\nPlease review the updated security policy and ensure compliance. Stay safe out there!\n\nThanks,\nKaren",
	},
	{
		Id:          11,
		Sender:      "larry.lane@company.com",
		Description: "Office Renovation",
		Body:        "Team,\n\nHeads up: The office renovation starts next week. Expect some noise and redirect any construction questions my way!\n\nBest,\nLarry",
	},
}
