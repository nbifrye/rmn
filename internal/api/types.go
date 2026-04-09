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
	ID        int    `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Mail      string `json:"mail"`
	CreatedOn string `json:"created_on"`
}
