package serializer

import "github.com/phatnt199/go-infra/pkg/core/metadata"

type MetadataSerializer interface {
	Serialize(meta metadata.Metadata) ([]byte, error)
	Deserialize(bytes []byte) (metadata.Metadata, error)
}
