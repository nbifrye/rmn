package api

type IdName struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Issue represents a Redmine issue.
type Issue struct {
	ID          int     `json:"id"`
	Project     IdName  `json:"project"`
	Tracker     IdName  `json:"tracker"`
	Status      IdName  `json:"status"`
	Priority    IdName  `json:"priority"`
	Author      IdName  `json:"author"`
	AssignedTo  *IdName `json:"assigned_to,omitempty"`
	Subject     string  `json:"subject"`
	Description string  `json:"description"`
	DoneRatio   int     `json:"done_ratio"`
	CreatedOn   string  `json:"created_on"`
	UpdatedOn   string  `json:"updated_on"`
}

// IssueCreateParams holds parameters for creating an issue.
type IssueCreateParams struct {
	ProjectID   int    `json:"project_id"`
	TrackerID   int    `json:"tracker_id,omitempty"`
	StatusID    int    `json:"status_id,omitempty"`
	PriorityID  int    `json:"priority_id,omitempty"`
	Subject     string `json:"subject"`
	Description string `json:"description,omitempty"`
	AssignedToID int   `json:"assigned_to_id,omitempty"`
}

// IssueUpdateParams holds parameters for updating an issue.
type IssueUpdateParams struct {
	TrackerID    int    `json:"tracker_id,omitempty"`
	StatusID     int    `json:"status_id,omitempty"`
	PriorityID   int    `json:"priority_id,omitempty"`
	Subject      string `json:"subject,omitempty"`
	Description  string `json:"description,omitempty"`
	AssignedToID int    `json:"assigned_to_id,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

// IssueListParams holds parameters for listing issues.
type IssueListParams struct {
	ProjectID    int    `json:"project_id,omitempty"`
	StatusID     string `json:"status_id,omitempty"`
	AssignedToID string `json:"assigned_to_id,omitempty"`
	TrackerID    int    `json:"tracker_id,omitempty"`
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
