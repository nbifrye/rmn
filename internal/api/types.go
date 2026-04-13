package api

type IdName struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IssueParent represents a parent issue reference.
type IssueParent struct {
	ID int `json:"id"`
}

// CustomField represents a Redmine custom field value.
type CustomField struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // string or []string for multi-value fields
}

// Journal represents a change history entry on an issue.
type Journal struct {
	ID        int             `json:"id"`
	User      IdName          `json:"user"`
	Notes     string          `json:"notes"`
	CreatedOn string          `json:"created_on"`
	Details   []JournalDetail `json:"details,omitempty"`
}

// JournalDetail represents a single field change within a journal entry.
type JournalDetail struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// Attachment represents a file attached to an issue.
type Attachment struct {
	ID          int    `json:"id"`
	Filename    string `json:"filename"`
	Filesize    int    `json:"filesize"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
	ContentURL  string `json:"content_url"`
	Author      IdName `json:"author"`
	CreatedOn   string `json:"created_on"`
}

// Relation represents a relationship between two issues.
type Relation struct {
	ID           int    `json:"id"`
	IssueID      int    `json:"issue_id"`
	IssueToID    int    `json:"issue_to_id"`
	RelationType string `json:"relation_type"`
	Delay        *int   `json:"delay"`
}

// IssueChild represents a child issue in the children list.
type IssueChild struct {
	ID      int    `json:"id"`
	Tracker IdName `json:"tracker"`
	Subject string `json:"subject"`
}

// Issue represents a Redmine issue.
type Issue struct {
	ID                  int           `json:"id"`
	Project             IdName        `json:"project"`
	Tracker             IdName        `json:"tracker"`
	Status              IdName        `json:"status"`
	Priority            IdName        `json:"priority"`
	Author              IdName        `json:"author"`
	AssignedTo          *IdName       `json:"assigned_to,omitempty"`
	Category            *IdName       `json:"category,omitempty"`
	FixedVersion        *IdName       `json:"fixed_version,omitempty"`
	Parent              *IssueParent  `json:"parent,omitempty"`
	Subject             string        `json:"subject"`
	Description         string        `json:"description"`
	StartDate           *string       `json:"start_date"`
	DueDate             *string       `json:"due_date"`
	DoneRatio           int           `json:"done_ratio"`
	IsPrivate           bool          `json:"is_private"`
	EstimatedHours      *float64      `json:"estimated_hours"`
	TotalEstimatedHours *float64      `json:"total_estimated_hours"`
	SpentHours          *float64      `json:"spent_hours"`
	TotalSpentHours     *float64      `json:"total_spent_hours"`
	CustomFields        []CustomField `json:"custom_fields,omitempty"`
	CreatedOn           string        `json:"created_on"`
	UpdatedOn           string        `json:"updated_on"`
	ClosedOn            *string       `json:"closed_on"`
	Journals            []Journal     `json:"journals,omitempty"`
	Attachments         []Attachment  `json:"attachments,omitempty"`
	Relations           []Relation    `json:"relations,omitempty"`
	Children            []IssueChild  `json:"children,omitempty"`
	Watchers            []IdName      `json:"watchers,omitempty"`
}

// IssueCreateParams holds parameters for creating an issue.
// ProjectID accepts both numeric IDs and string identifiers (e.g. "my-project").
type IssueCreateParams struct {
	ProjectID      interface{} `json:"project_id"`
	TrackerID      int         `json:"tracker_id,omitempty"`
	StatusID       int         `json:"status_id,omitempty"`
	PriorityID     int         `json:"priority_id,omitempty"`
	Subject        string      `json:"subject"`
	Description    string      `json:"description,omitempty"`
	AssignedToID   int         `json:"assigned_to_id,omitempty"`
	CategoryID     int         `json:"category_id,omitempty"`
	FixedVersionID int         `json:"fixed_version_id,omitempty"`
	ParentIssueID  int         `json:"parent_issue_id,omitempty"`
	StartDate      string      `json:"start_date,omitempty"`
	DueDate        string      `json:"due_date,omitempty"`
	EstimatedHours float64     `json:"estimated_hours,omitempty"`
	DoneRatio      int         `json:"done_ratio,omitempty"`
	IsPrivate      bool        `json:"is_private,omitempty"`
}

// IssueUpdateParams holds parameters for updating an issue.
// Pointer fields distinguish "not provided" (nil) from "set to zero value".
type IssueUpdateParams struct {
	TrackerID      *int     `json:"tracker_id,omitempty"`
	StatusID       *int     `json:"status_id,omitempty"`
	PriorityID     *int     `json:"priority_id,omitempty"`
	Subject        *string  `json:"subject,omitempty"`
	Description    *string  `json:"description,omitempty"`
	AssignedToID   *int     `json:"assigned_to_id,omitempty"`
	CategoryID     *int     `json:"category_id,omitempty"`
	FixedVersionID *int     `json:"fixed_version_id,omitempty"`
	ParentIssueID  *int     `json:"parent_issue_id,omitempty"`
	StartDate      *string  `json:"start_date,omitempty"`
	DueDate        *string  `json:"due_date,omitempty"`
	EstimatedHours *float64 `json:"estimated_hours,omitempty"`
	DoneRatio      *int     `json:"done_ratio,omitempty"`
	IsPrivate      *bool    `json:"is_private,omitempty"`
	Notes          string   `json:"notes,omitempty"`
	PrivateNotes   bool     `json:"private_notes,omitempty"`
}

// IntPtr returns a pointer to the given int value.
func IntPtr(v int) *int { return &v }

// StringPtr returns a pointer to the given string value.
func StringPtr(v string) *string { return &v }

// Float64Ptr returns a pointer to the given float64 value.
func Float64Ptr(v float64) *float64 { return &v }

// BoolPtr returns a pointer to the given bool value.
func BoolPtr(v bool) *bool { return &v }

// IssueListParams holds parameters for listing issues.
// ProjectID accepts both numeric IDs and string identifiers (e.g. "my-project").
type IssueListParams struct {
	ProjectID    string `json:"project_id,omitempty"`
	StatusID     string `json:"status_id,omitempty"`
	AssignedToID string `json:"assigned_to_id,omitempty"`
	TrackerID    int    `json:"tracker_id,omitempty"`
	Sort         string `json:"sort,omitempty"`
	Offset       int    `json:"offset,omitempty"`
	Limit        int    `json:"limit,omitempty"`
}

// PaginatedResponse wraps paginated API responses.
type PaginatedResponse struct {
	TotalCount int `json:"total_count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
}

// User represents a Redmine user.
type User struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Mail        string `json:"mail"`
	Admin       bool   `json:"admin,omitempty"`
	Status      int    `json:"status,omitempty"`
	LastLoginOn string `json:"last_login_on,omitempty"`
	CreatedOn   string `json:"created_on"`
}

// UserListParams holds parameters for listing users.
type UserListParams struct {
	Status  int    // 1=active, 2=registered, 3=locked (0 means no filter)
	Name    string // filter by name/login
	GroupID int
	Offset  int
	Limit   int
}

// IssueStatus represents a Redmine issue status.
type IssueStatus struct {
	ID       int  `json:"id"`
	Name     string `json:"name"`
	IsClosed bool   `json:"is_closed"`
}

// Project represents a Redmine project.
type Project struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Identifier  string   `json:"identifier"`
	Description string   `json:"description"`
	Homepage    string   `json:"homepage,omitempty"`
	Status      int      `json:"status"`
	IsPublic    bool     `json:"is_public"`
	Parent      *IdName  `json:"parent,omitempty"`
	CreatedOn   string   `json:"created_on"`
	UpdatedOn   string   `json:"updated_on"`
	Trackers    []IdName `json:"trackers,omitempty"`
}

// ProjectCreateParams holds parameters for creating a project.
type ProjectCreateParams struct {
	Name           string `json:"name"`
	Identifier     string `json:"identifier"`
	Description    string `json:"description,omitempty"`
	Homepage       string `json:"homepage,omitempty"`
	IsPublic       bool   `json:"is_public,omitempty"`
	ParentID       int    `json:"parent_id,omitempty"`
	InheritMembers bool   `json:"inherit_members,omitempty"`
}

// ProjectUpdateParams holds parameters for updating a project.
// Pointer fields distinguish "not provided" (nil) from "set to zero value".
type ProjectUpdateParams struct {
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	Homepage       *string `json:"homepage,omitempty"`
	IsPublic       *bool   `json:"is_public,omitempty"`
	ParentID       *int    `json:"parent_id,omitempty"`
	InheritMembers *bool   `json:"inherit_members,omitempty"`
}

// ProjectListParams holds parameters for listing projects.
type ProjectListParams struct {
	Status string // active, closed, archived
	Offset int
	Limit  int
}

// Version represents a Redmine project version.
type Version struct {
	ID            int     `json:"id"`
	Project       IdName  `json:"project"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Status        string  `json:"status"`
	DueDate       *string `json:"due_date"`
	Sharing       string  `json:"sharing"`
	WikiPageTitle string  `json:"wiki_page_title,omitempty"`
	CreatedOn     string  `json:"created_on"`
	UpdatedOn     string  `json:"updated_on"`
}

// VersionCreateParams holds parameters for creating a version.
type VersionCreateParams struct {
	Name          string `json:"name"`
	Status        string `json:"status,omitempty"`
	Sharing       string `json:"sharing,omitempty"`
	DueDate       string `json:"due_date,omitempty"`
	Description   string `json:"description,omitempty"`
	WikiPageTitle string `json:"wiki_page_title,omitempty"`
}

// VersionUpdateParams holds parameters for updating a version.
type VersionUpdateParams struct {
	Name          *string `json:"name,omitempty"`
	Status        *string `json:"status,omitempty"`
	Sharing       *string `json:"sharing,omitempty"`
	DueDate       *string `json:"due_date,omitempty"`
	Description   *string `json:"description,omitempty"`
	WikiPageTitle *string `json:"wiki_page_title,omitempty"`
}

// TimeEntry represents a Redmine time entry.
type TimeEntry struct {
	ID        int     `json:"id"`
	Project   IdName  `json:"project"`
	Issue     *IdName `json:"issue,omitempty"`
	User      IdName  `json:"user"`
	Activity  IdName  `json:"activity"`
	Hours     float64 `json:"hours"`
	Comments  string  `json:"comments"`
	SpentOn   string  `json:"spent_on"`
	CreatedOn string  `json:"created_on"`
	UpdatedOn string  `json:"updated_on"`
}

// TimeEntryCreateParams holds parameters for creating a time entry.
type TimeEntryCreateParams struct {
	IssueID    int     `json:"issue_id,omitempty"`
	ProjectID  string  `json:"project_id,omitempty"`
	SpentOn    string  `json:"spent_on,omitempty"`
	Hours      float64 `json:"hours"`
	ActivityID int     `json:"activity_id,omitempty"`
	Comments   string  `json:"comments,omitempty"`
}

// TimeEntryUpdateParams holds parameters for updating a time entry.
type TimeEntryUpdateParams struct {
	IssueID    *int     `json:"issue_id,omitempty"`
	ProjectID  *string  `json:"project_id,omitempty"`
	SpentOn    *string  `json:"spent_on,omitempty"`
	Hours      *float64 `json:"hours,omitempty"`
	ActivityID *int     `json:"activity_id,omitempty"`
	Comments   *string  `json:"comments,omitempty"`
}

// TimeEntryListParams holds parameters for listing time entries.
type TimeEntryListParams struct {
	ProjectID  string
	IssueID    int
	UserID     int
	SpentOn    string
	From       string
	To         string
	ActivityID int
	Offset     int
	Limit      int
}

// Membership represents a Redmine project membership.
type Membership struct {
	ID      int      `json:"id"`
	Project IdName   `json:"project"`
	User    *IdName  `json:"user,omitempty"`
	Group   *IdName  `json:"group,omitempty"`
	Roles   []IdName `json:"roles"`
}

// MembershipCreateParams holds parameters for creating a membership.
type MembershipCreateParams struct {
	UserID  int   `json:"user_id"`
	RoleIDs []int `json:"role_ids"`
}

// MembershipUpdateParams holds parameters for updating a membership.
type MembershipUpdateParams struct {
	RoleIDs []int `json:"role_ids"`
}

// MembershipListParams holds parameters for listing memberships.
type MembershipListParams struct {
	Offset int
	Limit  int
}

// WikiPage represents a Redmine wiki page in index listings.
type WikiPage struct {
	Title     string  `json:"title"`
	Parent    *IdName `json:"parent,omitempty"`
	Version   int     `json:"version"`
	CreatedOn string  `json:"created_on"`
	UpdatedOn string  `json:"updated_on"`
}

// WikiPageDetail represents a full Redmine wiki page with content.
type WikiPageDetail struct {
	Title     string  `json:"title"`
	Parent    *IdName `json:"parent,omitempty"`
	Text      string  `json:"text"`
	Version   int     `json:"version"`
	Author    IdName  `json:"author"`
	Comments  string  `json:"comments"`
	CreatedOn string  `json:"created_on"`
	UpdatedOn string  `json:"updated_on"`
}

// WikiPageCreateParams holds parameters for creating a wiki page.
type WikiPageCreateParams struct {
	Text     string `json:"text"`
	Comments string `json:"comments,omitempty"`
}

// WikiPageUpdateParams holds parameters for updating a wiki page.
type WikiPageUpdateParams struct {
	Text     *string `json:"text,omitempty"`
	Comments *string `json:"comments,omitempty"`
	Version  *int    `json:"version,omitempty"`
}
