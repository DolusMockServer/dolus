package loader

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"

	"github.com/DolusMockServer/dolus/pkg/logger"
)

// TODO: (watch for cue package manager)
const dolusExpectationsHomeFolder = "cue/github.com/DolusMockServer/dolus/cue-expectations"

type (
	CueExpectationLoadType []cue.Value
)

type CueExpectationLoader struct {
	cueExpectationsFiles           []string
	cueDolusExpectationsRootModule string
}

var _ Loader[CueExpectationLoadType] = &CueExpectationLoader{}

func NewCueExpectationLoader(cueExpectationsFiles []string) *CueExpectationLoader {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Log.Fatalf("failed to user home dir: %s", err.Error())
	}
	return &CueExpectationLoader{
		// TODO: allow this to be passed as an argument
		cueDolusExpectationsRootModule: fmt.Sprintf("%s/%s", homeDir, dolusExpectationsHomeFolder),
		cueExpectationsFiles:           cueExpectationsFiles,
	}
}

func (cel *CueExpectationLoader) Load() (*CueExpectationLoadType, error) {
	ctx := cuecontext.New()
	var cueValues []cue.Value
	logger.Log.Debugf(
		"Loading expectations from cue root module: %s",
		cel.cueDolusExpectationsRootModule,
	)

	if len(cel.cueExpectationsFiles) == 0 {
		return (*CueExpectationLoadType)(&cueValues), nil
	}

	bis := load.Instances(cel.cueExpectationsFiles, &load.Config{
		ModuleRoot: cel.cueDolusExpectationsRootModule,
	})

	for _, bi := range bis {
		if bi.Err != nil {
			logger.Log.Error("Error during load: ", bi.Err)
			continue
		}
		value := ctx.BuildInstance(bi)

		if value.Err() != nil {
			logger.Log.Error("Error during load: ", value.Err())
			continue
		}

		err := value.Validate()
		if err != nil {
			logger.Log.Error("Error during validation:", err)
			continue
		}

		cueValues = append(cueValues, value)
	}

	return (*CueExpectationLoadType)(&cueValues), nil
}
