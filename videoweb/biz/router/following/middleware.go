// Code generated by hertz generator.

package following

import (
	"github.com/cloudwego/hertz/pkg/app"
	"videoweb/middleware"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _followingMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _listMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{middleware.AuthMiddleware()}
}