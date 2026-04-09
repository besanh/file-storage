package biz

type CheckPermissionRequest struct {
	SubjectType string
	SubjectID   string
	Relation    string
	ObjectType  string
	ObjectID    string
}

type CheckPermissionResponse struct {
	Allowed bool
}

type WriteRelationshipRequest struct {
	ResourceType string
	ResourceID   string
	Relation     string
	SubjectType  string
	SubjectID    string
}

type WriteRelationshipResponse struct {
}

type DeleteRelationshipRequest struct {
	ResourceType string
	ResourceID   string
	Relation     string
	SubjectType  string
	SubjectID    string
}

type DeleteRelationshipResponse struct {
}

type SwapRelationshipRequest struct {
	ResourceType string
	ResourceID   string
	SubjectType  string
	SubjectID    string
	OldRelation  string
	NewRelation  string
}

type SwapRelationshipResponse struct {
}
