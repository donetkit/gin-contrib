package minio

import "time"

// UploadInfo contains information about the
// newly uploaded or copied object.
type UploadInfo struct {
	Bucket       string
	Key          string
	ETag         string
	Size         int64
	LastModified time.Time
}

// ObjectInfo container for object metadata.
type ObjectInfo struct {
	// An ETag is optionally set to md5sum of an object.  In case of multipart objects,
	// ETag is of the form MD5SUM-N where MD5SUM is md5sum of all individual md5sums of
	// each parts concatenated into one string.
	ETag string `json:"etag"`

	Key          string    `json:"name"`         // Name of the object
	LastModified time.Time `json:"lastModified"` // Date and time the object was last modified.
	Size         int64     `json:"size"`         // Size in bytes of the object.
	ContentType  string    `json:"contentType"`  // A standard MIME type describing the format of the object data.

	// x-amz-tagging values in their k/v values.
	UserTags map[string]string `json:"userTags"`

	// x-amz-tagging-count value
	UserTagCount int

	// The class of storage used to store the object.
	StorageClass string `json:"storageClass"`

	// Versioning related information
	IsLatest       bool
	IsDeleteMarker bool
	VersionID      string `xml:"VersionId"`
}
