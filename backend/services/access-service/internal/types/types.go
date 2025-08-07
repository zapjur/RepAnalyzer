package types

type GrpcDBServiceResponse struct {
	Owned     bool
	Message   string
	ObjectKey string
	Bucket    string
}
