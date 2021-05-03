/*
Package client implements the underlying API calling logic.

Initialize the `Client` and pass it to other API services.

Example:

	cli := New(apiKey, username, sandbox)
	// You can pass this to the service methods
	rep, err := SendMessage(cli, req)
*/
package client
