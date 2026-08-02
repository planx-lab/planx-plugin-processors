package processor

import (
	"context"

	"github.com/planx-lab/planx-sdk-go/sdk"
)

// Passthrough returns the batch unchanged (1:1). No config is required.
type Passthrough struct{}

// NewPassthrough builds the passthrough processor.
func NewPassthrough() sdk.ProcessorSPI { return &Passthrough{} }

// Init accepts any config (including empty/nil) and is a no-op.
func (p *Passthrough) Init(_ context.Context, _ []byte) error { return nil }

// Process returns the batch exactly as received.
func (p *Passthrough) Process(b sdk.Batch) (sdk.Batch, error) { return b, nil }

func (p *Passthrough) Close() error { return nil }
