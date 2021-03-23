package actions

import (

  "fmt"
  "github.com/aeolus3000/lendo-sdk/banking"
  "github.com/aeolus3000/lendo-sdk/messaging"
  "lendo_service/middleware"
  "net/http"
  "github.com/gobuffalo/buffalo"
  "github.com/gobuffalo/pop/v5"
  "github.com/gobuffalo/x/responder"
  "lendo_service/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Application)
// DB Table: Plural (applications)
// Resource: Plural (Applications)
// Path: Plural (/applications)
// View Template Folder: Plural (/templates/applications/)

// ApplicationsResource is the resource for the Application model
type ApplicationsResource struct{
  buffalo.Resource
}

// List gets all Applications. This function is mapped to the path
// GET /applications
func (v ApplicationsResource) List(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  applications := &models.Applications{}

  // Paginate results. Params "page" and "per_page" control pagination.
  // Default values are "page=1" and "per_page=20".
  q := tx.PaginateFromParams(c.Params())

  fmt.Println("Application status from url: " + c.Param("with_status"))

  // Retrieve all Applications from the DB
  if err := q.All(applications); err != nil {
    return err
  }

  return responder.Wants("html", func (c buffalo.Context) error {
    // Add the paginator to the context so it can be used in the template.
    c.Set("pagination", q.Paginator)

    c.Set("applications", applications)
    return c.Render(http.StatusOK, r.HTML("/applications/index.plush.html"))
  }).Wants("json", func (c buffalo.Context) error {
    return c.Render(200, r.JSON(applications))
  }).Wants("xml", func (c buffalo.Context) error {
    return c.Render(200, r.XML(applications))
  }).Respond(c)
}

// Show gets the data for one Application. This function is mapped to
// the path GET /applications/{application_id}
func (v ApplicationsResource) Show(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty Application
  application := &models.Application{}

  // To find the Application the parameter application_id is used.
  if err := tx.Find(application, c.Param("application_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  return responder.Wants("html", func (c buffalo.Context) error {
    c.Set("application", application)

    return c.Render(http.StatusOK, r.HTML("/applications/show.plush.html"))
  }).Wants("json", func (c buffalo.Context) error {
    return c.Render(200, r.JSON(application))
  }).Wants("xml", func (c buffalo.Context) error {
    return c.Render(200, r.XML(application))
  }).Respond(c)
}

// New renders the form for creating a new Application.
// This function is mapped to the path GET /applications/new
func (v ApplicationsResource) New(c buffalo.Context) error {
  c.Set("application", &models.Application{})

  return c.Render(http.StatusOK, r.HTML("/applications/new.plush.html"))
}
// Create adds a Application to the DB. This function is mapped to the
// path POST /applications
func (v ApplicationsResource) Create(c buffalo.Context) error {
  // Allocate an empty Application
  application := &models.Application{}

  // Bind application to the html form elements
  if err := c.Bind(application); err != nil {
    return err
  }

  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  pub, ok := c.Value(middleware.CtxPublisher).(messaging.Publisher)

  bankingApplication := banking.Application{
    Id:        application.ID.String(),
    FirstName: application.FirstName,
    LastName:  application.LastName.String,
    Status:    application.Status.String,
    JobId:     application.JobID.String,
  }
  bytesBuffer, serializeError := banking.SerializeFromApplication(&bankingApplication)
  if serializeError != nil {
    return serializeError
  }
  publishErr := pub.Publish(bytesBuffer)
  if publishErr != nil {
    return publishErr
  }

  // Validate the data from the html form
  verrs, err := tx.ValidateAndCreate(application)
  if err != nil {
    return err
  }

  if verrs.HasAny() {
    return responder.Wants("html", func (c buffalo.Context) error {
      // Make the errors available inside the html template
      c.Set("errors", verrs)

      // Render again the new.html template that the user can
      // correct the input.
      c.Set("application", application)

      return c.Render(http.StatusUnprocessableEntity, r.HTML("/applications/new.plush.html"))
    }).Wants("json", func (c buffalo.Context) error {
      return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
    }).Wants("xml", func (c buffalo.Context) error {
      return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
    }).Respond(c)
  }

  return responder.Wants("html", func (c buffalo.Context) error {
    // If there are no errors set a success message
    c.Flash().Add("success", T.Translate(c, "application.created.success"))

    // and redirect to the show page
    return c.Redirect(http.StatusSeeOther, "/applications/%v", application.ID)
  }).Wants("json", func (c buffalo.Context) error {
    return c.Render(http.StatusCreated, r.JSON(application))
  }).Wants("xml", func (c buffalo.Context) error {
    return c.Render(http.StatusCreated, r.XML(application))
  }).Respond(c)
}

// Edit renders a edit form for a Application. This function is
// mapped to the path GET /applications/{application_id}/edit
func (v ApplicationsResource) Edit(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty Application
  application := &models.Application{}

  if err := tx.Find(application, c.Param("application_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  c.Set("application", application)
  return c.Render(http.StatusOK, r.HTML("/applications/edit.plush.html"))
}
// Update changes a Application in the DB. This function is mapped to
// the path PUT /applications/{application_id}
func (v ApplicationsResource) Update(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty Application
  application := &models.Application{}

  if err := tx.Find(application, c.Param("application_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  // Bind Application to the html form elements
  if err := c.Bind(application); err != nil {
    return err
  }

  verrs, err := tx.ValidateAndUpdate(application)
  if err != nil {
    return err
  }

  if verrs.HasAny() {
    return responder.Wants("html", func (c buffalo.Context) error {
      // Make the errors available inside the html template
      c.Set("errors", verrs)

      // Render again the edit.html template that the user can
      // correct the input.
      c.Set("application", application)

      return c.Render(http.StatusUnprocessableEntity, r.HTML("/applications/edit.plush.html"))
    }).Wants("json", func (c buffalo.Context) error {
      return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
    }).Wants("xml", func (c buffalo.Context) error {
      return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
    }).Respond(c)
  }

  return responder.Wants("html", func (c buffalo.Context) error {
    // If there are no errors set a success message
    c.Flash().Add("success", T.Translate(c, "application.updated.success"))

    // and redirect to the show page
    return c.Redirect(http.StatusSeeOther, "/applications/%v", application.ID)
  }).Wants("json", func (c buffalo.Context) error {
    return c.Render(http.StatusOK, r.JSON(application))
  }).Wants("xml", func (c buffalo.Context) error {
    return c.Render(http.StatusOK, r.XML(application))
  }).Respond(c)
}

// Destroy deletes a Application from the DB. This function is mapped
// to the path DELETE /applications/{application_id}
func (v ApplicationsResource) Destroy(c buffalo.Context) error {
  // Get the DB connection from the context
  tx, ok := c.Value("tx").(*pop.Connection)
  if !ok {
    return fmt.Errorf("no transaction found")
  }

  // Allocate an empty Application
  application := &models.Application{}

  // To find the Application the parameter application_id is used.
  if err := tx.Find(application, c.Param("application_id")); err != nil {
    return c.Error(http.StatusNotFound, err)
  }

  if err := tx.Destroy(application); err != nil {
    return err
  }

  return responder.Wants("html", func (c buffalo.Context) error {
    // If there are no errors set a flash message
    c.Flash().Add("success", T.Translate(c, "application.destroyed.success"))

    // Redirect to the index page
    return c.Redirect(http.StatusSeeOther, "/applications")
  }).Wants("json", func (c buffalo.Context) error {
    return c.Render(http.StatusOK, r.JSON(application))
  }).Wants("xml", func (c buffalo.Context) error {
    return c.Render(http.StatusOK, r.XML(application))
  }).Respond(c)
}