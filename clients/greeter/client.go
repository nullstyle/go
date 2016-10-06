package greeter

// Close shuts down the client, closing the underlying grpc connection
func (c *Client) Close() error {
	return c.conn.Close()
}
