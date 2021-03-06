package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAppsCommand(t *testing.T) {
	cmd := Apps()
	require.NotNil(t, cmd)
	assertCommandNames(t, cmd,
		"create",
		"get",
		"list",
		"update",
		"delete",
		"create-deployment",
		"get-deployment",
		"list-deployments",
		"logs",
	)
}

var testAppSpec = godo.AppSpec{
	Name: "test",
	Services: []godo.AppServiceSpec{
		{
			Name: "service",
			GitHub: godo.GitHubSourceSpec{
				Repo:   "digitalocean/doctl",
				Branch: "master",
			},
		},
	},
}

func TestRunAppsCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		specFile, err := ioutil.TempFile("", "spec")
		require.NoError(t, err)
		defer func() {
			os.Remove(specFile.Name())
			specFile.Close()
		}()

		err = json.NewEncoder(specFile).Encode(&testAppSpec)
		require.NoError(t, err)

		app := &godo.App{
			ID:        uuid.New().String(),
			Spec:      &testAppSpec,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		createReq := &godo.AppCreateRequest{
			Spec: &testAppSpec,
		}

		tm.apps.EXPECT().Create(createReq).Times(1).Return(app, nil)

		config.Doit.Set(config.NS, doctl.ArgAppSpec, specFile.Name())

		err = RunAppsCreate(config)
		require.NoError(t, err)
	})
}

func TestRunAppsGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		app := &godo.App{
			ID:        uuid.New().String(),
			Spec:      &testAppSpec,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tm.apps.EXPECT().Get(app.ID).Times(1).Return(app, nil)

		config.Args = append(config.Args, app.ID)

		err := RunAppsGet(config)
		require.NoError(t, err)
	})
}

func TestRunAppsList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		apps := []*godo.App{{
			ID:        uuid.New().String(),
			Spec:      &testAppSpec,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}}

		tm.apps.EXPECT().List().Times(1).Return(apps, nil)

		err := RunAppsList(config)
		require.NoError(t, err)
	})
}

func TestRunAppsUpdate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		specFile, err := ioutil.TempFile("", "spec")
		require.NoError(t, err)
		defer func() {
			os.Remove(specFile.Name())
			specFile.Close()
		}()

		err = json.NewEncoder(specFile).Encode(&testAppSpec)
		require.NoError(t, err)

		app := &godo.App{
			ID:        uuid.New().String(),
			Spec:      &testAppSpec,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		updateReq := &godo.AppUpdateRequest{
			Spec: &testAppSpec,
		}

		tm.apps.EXPECT().Update(app.ID, updateReq).Times(1).Return(app, nil)

		config.Args = append(config.Args, app.ID)
		config.Doit.Set(config.NS, doctl.ArgAppSpec, specFile.Name())

		err = RunAppsUpdate(config)
		require.NoError(t, err)
	})
}

func TestRunAppsDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		app := &godo.App{
			ID:        uuid.New().String(),
			Spec:      &testAppSpec,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tm.apps.EXPECT().Delete(app.ID).Times(1).Return(nil)

		config.Args = append(config.Args, app.ID)

		err := RunAppsDelete(config)
		require.NoError(t, err)
	})
}

func TestRunAppsCreateDeployment(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		appID := uuid.New().String()
		deployment := &godo.Deployment{
			ID:   uuid.New().String(),
			Spec: &testAppSpec,
			Services: []*godo.DeploymentService{{
				Name:             "service",
				SourceCommitHash: "commit",
			}},
			Cause: "Manual",
			Progress: &godo.DeploymentProgress{
				PendingSteps: 1,
				RunningSteps: 0,
				SuccessSteps: 0,
				ErrorSteps:   0,
				TotalSteps:   1,

				Steps: []*godo.DeploymentProgressStep{{
					Name:      "name",
					Status:    "pending",
					Attempts:  0,
					StartedAt: time.Now(),
				}},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tm.apps.EXPECT().CreateDeployment(appID).Times(1).Return(deployment, nil)

		config.Args = append(config.Args, appID)

		err := RunAppsCreateDeployment(config)
		require.NoError(t, err)
	})
}

func TestRunAppsGetDeployment(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		appID := uuid.New().String()
		deployment := &godo.Deployment{
			ID:   uuid.New().String(),
			Spec: &testAppSpec,
			Services: []*godo.DeploymentService{{
				Name:             "service",
				SourceCommitHash: "commit",
			}},
			Cause: "Manual",
			Progress: &godo.DeploymentProgress{
				PendingSteps: 1,
				RunningSteps: 0,
				SuccessSteps: 0,
				ErrorSteps:   0,
				TotalSteps:   1,

				Steps: []*godo.DeploymentProgressStep{{
					Name:      "name",
					Status:    "pending",
					Attempts:  0,
					StartedAt: time.Now(),
				}},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		tm.apps.EXPECT().GetDeployment(appID, deployment.ID).Times(1).Return(deployment, nil)

		config.Args = append(config.Args, appID, deployment.ID)

		err := RunAppsGetDeployment(config)
		require.NoError(t, err)
	})
}

func TestRunAppsListDeployments(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		appID := uuid.New().String()
		deployments := []*godo.Deployment{{
			ID:   uuid.New().String(),
			Spec: &testAppSpec,
			Services: []*godo.DeploymentService{{
				Name:             "service",
				SourceCommitHash: "commit",
			}},
			Cause: "Manual",
			Progress: &godo.DeploymentProgress{
				PendingSteps: 1,
				RunningSteps: 0,
				SuccessSteps: 0,
				ErrorSteps:   0,
				TotalSteps:   1,

				Steps: []*godo.DeploymentProgressStep{{
					Name:      "name",
					Status:    "pending",
					Attempts:  0,
					StartedAt: time.Now(),
				}},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}}

		tm.apps.EXPECT().ListDeployments(appID).Times(1).Return(deployments, nil)

		config.Args = append(config.Args, appID)

		err := RunAppsListDeployments(config)
		require.NoError(t, err)
	})
}

func TestRunAppsGetLogs(t *testing.T) {
	appID := uuid.New().String()
	deploymentID := uuid.New().String()
	component := "service"

	types := map[string]godo.AppLogType{
		"build":  godo.AppLogTypeBuild,
		"deploy": godo.AppLogTypeDeploy,
		"run":    godo.AppLogTypeRun,
	}

	for typeStr, logType := range types {
		withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
			tm.apps.EXPECT().GetLogs(appID, deploymentID, component, logType).Times(1).Return(&godo.AppLogs{LiveURL: "https://proxy-apps-prod-ams3-001.ondigitalocean.app/?token=..."}, nil)

			config.Args = append(config.Args, appID, deploymentID, component)
			config.Doit.Set(config.NS, doctl.ArgAppLogType, typeStr)

			err := RunAppsGetLogs(config)
			require.NoError(t, err)
		})
	}
}
