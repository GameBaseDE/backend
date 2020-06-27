/*
 * GameBase Communication API
 *
 * This is the REST API used as an communication layer between backend and frontend.
 *
 * API version: 2.1.0
 * Contact: gamebase@gahr.dev
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Exception struct {

	// ID of container causing this exception
	Id string `json:"id,omitempty"`

	// Short Exception code
	Exception string `json:"exception"`

	// Detailed message of issue
	Details string `json:"details"`
}
