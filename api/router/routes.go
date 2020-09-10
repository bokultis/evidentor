package router

import (
	"github.com/bokultis/evidentor/api/student"
	"github.com/bokultis/evidentor/api/user"
)

var rootRoutes = RoutePrefix{
	"/",
	[]Route{

		Route{
			"UsersLogin",
			"POST",
			"/login",
			user.LoginHandler,
			false,
		},
		Route{
			"UsersLogout",
			"GET",
			"/logout",
			user.LogoutHandler,
			true,
		},
	},
}

var userRoutes = RoutePrefix{
	"/users",
	[]Route{
		Route{
			"ListHandler",
			"GET",
			"",
			user.UsersIndexHandler,
			true,
		},
		Route{
			"ShowHandler",
			"GET",
			"/{userId}",
			user.UsersShowHandler,
			true,
		},
		Route{
			"CreateHandler",
			"POST",
			"",
			user.UsersCreateHandler,
			true,
		},

		Route{
			"DeleteHandler",
			"DELETE",
			"/{userId}",
			user.UsersDeleteHandler,
			true,
		},
		Route{
			"UpdateHandler",
			"PUT",
			"/{userId}",
			user.UsersUpdateHandler,
			true,
		},
	},
}

var studentRoutes = RoutePrefix{
	"/students",
	[]Route{
		Route{
			"ListHandler",
			"GET",
			"",
			student.StudentsIndexHandler,
			true,
		},
		Route{
			"ShowHandler",
			"GET",
			"/{studentId}",
			student.StudentsShowHandler,
			true,
		},

		Route{
			"CreateHandler",
			"POST",
			"",
			student.StudentsCreateHandler,
			true,
		},

		Route{
			"DeleteHandler",
			"DELETE",
			"/{studentId}",
			student.StudentsDeleteHandler,
			true,
		},
		Route{
			"UpdateHandler",
			"PUT",
			"/{studentId}",
			student.StudentsUpdateHandler,
			true,
		},
		Route{
			"ListGroupsHandler",
			"GET",
			"/groups/{studentId}",
			student.StudentsGroupsListHandler,
			true,
		},
		Route{
			"ListNotesHandler",
			"GET",
			"/notes/{studentId}",
			student.StudentsNotesListHandler,
			true,
		},
	},
}
