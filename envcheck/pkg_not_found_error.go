package envcheck

import "fmt"

func (err *PkgNotFoundError) Error() string {
	return fmt.Sprintf("cannot find pkg %s", err.Pkg)
}
