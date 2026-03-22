// Package migrator provides an in code database migration management it can be
// used to run migrations before starting a service, but better way would be to
// embed your database migrations and then have a separate job/mini-service that
// keeps it in sync, either way its up to the end user to do as they please
package migrator

// Migrator interface allows for generic migration application
type Migrator interface{}
