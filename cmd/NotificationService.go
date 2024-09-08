package cmd

import "sync"

type NotificationService interface {
	// send one time alert for validator
	SendValidatorAlertNotification(config *HalfLifeConfig, vm *ValidatorMonitor, stats ValidatorStats, alertNotification *ValidatorAlertNotification)

	// update (or create) realtime status for validator
	UpdateValidatorRealtimeStatus(statusFile string, config *HalfLifeConfig, status *HalfLifeStatus, vm *ValidatorMonitor, vmStatus *ValidatorStatus, stats ValidatorStats, writeConfigMutex *sync.Mutex)
}
