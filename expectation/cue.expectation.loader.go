package expectation

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/DolusMockServer/dolus/logger"
)

type (
	CueExpectationLoadType []cue.Value
)

type CueExpectationLoader struct {
	cueExpectationsFiles           []string
	cueDolusExpectationsRootModule string
}

const dolusExpectationsHomeFolder = "cue/github.com/MartinSimango/dolus-expectations"

var _ Loader[CueExpectationLoadType] = &CueExpectationLoader{}

func NewCueExpectationLoader(cueExpectationsFiles []string) *CueExpectationLoader {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Log.Fatalf("failed to user home dir: %s", err.Error())
	}
	return &CueExpectationLoader{
		cueDolusExpectationsRootModule: fmt.Sprintf("%s/%s", homeDir, dolusExpectationsHomeFolder),
		cueExpectationsFiles:           cueExpectationsFiles,
	}
}

func (cel *CueExpectationLoader) load() (*CueExpectationLoadType, error) {
	ctx := cuecontext.New()
	var cueValues []cue.Value
	logger.Log.Debugf("Loading expectations from cue root module: %s", cel.cueDolusExpectationsRootModule)
	bis := load.Instances(cel.cueExpectationsFiles, &load.Config{
		ModuleRoot: cel.cueDolusExpectationsRootModule,
	})

	for _, bi := range bis {
		// check for errors on the  instance
		// these are typically parsing errors
		if bi.Err != nil {
			logger.Log.Error("Error during load:", bi.Err)
			continue
		}
		value := ctx.BuildInstance(bi)

		if value.Err() != nil {
			logger.Log.Error("Error during load:", value.Err())
			continue
		}

		// Validate the value
		err := value.Validate()
		if err != nil {
			logger.Log.Error("Error during validation:", err)
			continue
		}

		cueValues = append(cueValues, value)
	}

	return (*CueExpectationLoadType)(&cueValues), nil
}
