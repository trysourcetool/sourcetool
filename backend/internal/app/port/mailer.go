package port

import "context"

type Mailer interface {
	Send(ctx context.Context, to []string, from, subject, body string) error
}
