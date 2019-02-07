package web

import (
	"html/template"
	"net/http"
)

type TemplateVarLanding struct {
	UserID  string
}

func HandleLanding(response http.ResponseWriter, request *http.Request) {
	us, err := globalSessions.Get(request, "session-key")
	if err != nil {
		MakeErrorResponse(response, 500, err.Error(), 0)
		return
	}

	tmlpStr, err := templates.FindString("templates/landing.html")
	if err != nil {
		MakeErrorResponse(response, 500, err.Error(), 0)
		return
	}

	templateVars := &TemplateVarLanding{}

	if us.Values["LoggedInUserID"] != nil {
		templateVars.UserID = us.Values["LoggedInUserID"].(string)
	}

	if err = us.Save(request, response); err != nil {
		logger.Errorf("Error saving session: %v", err)
		MakeErrorResponse(response, 500, err.Error(), 0)
		return
	}

	tmpl := template.New("landing template")
	tmpl = template.Must(tmpl.Parse(tmlpStr))
	tmpl.Execute(response, templateVars)
	return
}