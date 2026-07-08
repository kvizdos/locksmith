package events

import "context"

type contextKey struct{}

type ContextMetadata struct {
	RequestID string
	TraceID   string
	Source    string
	Values    map[string]string
}

func WithContextMetadata(ctx context.Context, metadata ContextMetadata) context.Context {
	metadata.Values = cloneValues(metadata.Values)
	return context.WithValue(ctx, contextKey{}, metadata)
}

func MetadataFromContext(ctx context.Context) ContextMetadata {
	metadata, _ := ctx.Value(contextKey{}).(ContextMetadata)
	metadata.Values = cloneValues(metadata.Values)
	return metadata
}

func EnrichEnvelope(ctx context.Context, event Envelope) Envelope {
	metadata := MetadataFromContext(ctx)
	if event.RequestID == "" {
		event.RequestID = metadata.RequestID
	}
	if event.TraceID == "" {
		event.TraceID = metadata.TraceID
	}
	if event.Source == "" {
		event.Source = metadata.Source
	}
	if event.Metadata == nil {
		event.Metadata = cloneValues(metadata.Values)
	} else {
		merged := cloneValues(metadata.Values)
		for key, value := range event.Metadata {
			merged[key] = value
		}
		event.Metadata = merged
	}
	return event
}

func cloneValues(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}
