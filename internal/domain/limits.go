package domain

import "fmt"

func (b ResourceBudget) Validate() error {
	if b.CPUSeconds < 0 || b.MemoryBytes < 0 || b.NetworkBytes < 0 {
		return fmt.Errorf("resource_budget_negative")
	}
	if b.MemoryBytes > 64<<30 {
		return fmt.Errorf("memory_budget_too_large")
	}
	return nil
}
func (b ResourceBudget) Fits(used ResourceBudget) bool {
	return used.CPUSeconds <= b.CPUSeconds && used.MemoryBytes <= b.MemoryBytes && used.NetworkBytes <= b.NetworkBytes
}
