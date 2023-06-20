package helpers

import (
	"errors"
	"fmt"

	"github.com/devopsarr/prowlarr-go/prowlarr"
)

// define constant for error management.
const (
	Create                            = "create"
	Read                              = "read"
	Update                            = "update"
	Delete                            = "delete"
	List                              = "list"
	ClientError                       = "Client Error"
	ResourceError                     = "Resource Error"
	DataSourceError                   = "Data Source Error"
	UnexpectedImportIdentifier        = "Unexpected Import Identifier"
	UnexpectedResourceConfigureType   = "Unexpected Resource Configure Type"
	UnexpectedDataSourceConfigureType = "Unexpected DataSource Configure Type"
)

var ErrDataNotFound = errors.New("data source not found")

func ErrDataNotFoundError(kind, field, search string) error {
	return fmt.Errorf("%w: no %s with %s '%s'", ErrDataNotFound, kind, field, search)
}

func ParseNotFoundError(kind, field, search string) string {
	return fmt.Sprintf("Unable to find %s, got error: data source not found: no %s with %s '%s'", kind, kind, field, search)
}

func ParseClientError(action, name string, err error) string {
	if e, ok := err.(*prowlarr.GenericOpenAPIError); ok {
		return fmt.Sprintf("Unable to %s %s, got error: %s\nDetails:\n%s", action, name, err, string(e.Body()))
	}

	return fmt.Sprintf("Unable to %s %s, got error: %s", action, name, err)
}
