package main

package main

import (
  "gorm.io/gorm"
)

type Product struct {
  gorm.Model
  Name string
  Email string
  Paswword string
  Age int
}
func main() {
	
}
