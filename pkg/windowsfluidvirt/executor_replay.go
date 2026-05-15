package windowsfluidvirt

import (
	"encoding/json"
	"os"
	"time"
)

type ExecutorReplayFixture struct {
	Name                   string                                   `json:"name"`
	GovernanceContract     WindowsFluidApplyGovernanceContract      `json:"governanceContract"`
	Revalidation           WindowsFluidPreApplyRevalidationContract `json:"revalidation"`
	Attestation            WindowsFluidPolicyAttestation            `json:"attestation"`
	KillSwitch             *WindowsFluidKillSwitch                  `json:"killSwitch,omitempty"`
	ExpectedExecutionPhase ExecutionPhase                           `json:"expectedExecutionPhase"`
	ExpectedBlockers       []string                                 `json:"expectedBlockers"`
}

func LoadExecutorReplayFixture(path string) (ExecutorReplayFixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ExecutorReplayFixture{}, err
	}
	var fixture ExecutorReplayFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return ExecutorReplayFixture{}, err
	}
	return fixture, nil
}

func EvaluateExecutorReplayFixture(path string, evaluationTime time.Time) (FutureApplyExecutorEvaluationResult, ExecutorReplayFixture, error) {
	fixture, err := LoadExecutorReplayFixture(path)
	if err != nil {
		return FutureApplyExecutorEvaluationResult{}, ExecutorReplayFixture{}, err
	}
	result := EvaluateWindowsFluidFutureApplyExecutor(FutureApplyExecutorEvaluationInput{
		GovernanceContract: fixture.GovernanceContract,
		Revalidation:       fixture.Revalidation,
		Attestation:        fixture.Attestation,
		KillSwitch:         fixture.KillSwitch,
		EvaluationTime:     evaluationTime,
	})
	return result, fixture, nil
}
