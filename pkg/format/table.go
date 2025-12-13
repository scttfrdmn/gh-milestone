package format

import (
	"fmt"
	"strings"
	"time"

	"github.com/scttfrdmn/gh-milestone-manager/pkg/api"
)

// FormatMilestoneTable formats milestones as a table
func FormatMilestoneTable(milestones []api.Milestone) string {
	if len(milestones) == 0 {
		return "No milestones found"
	}

	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("%-6s  %-40s  %-12s  %-15s\n", "NUMBER", "TITLE", "DUE DATE", "PROGRESS"))

	// Rows
	for _, m := range milestones {
		number := fmt.Sprintf("#%d", m.Number)
		title := truncate(m.Title, 40)
		dueDate := formatDueDate(m.DueOn)
		progress := formatProgress(m.OpenIssues, m.ClosedIssues)

		sb.WriteString(fmt.Sprintf("%-6s  %-40s  %-12s  %-15s\n", number, title, dueDate, progress))
	}

	return sb.String()
}

// FormatMilestoneDetail formats a milestone with full details
func FormatMilestoneDetail(m *api.Milestone) string {
	var sb strings.Builder

	// Title and state
	sb.WriteString(fmt.Sprintf("%s\n", m.Title))
	sb.WriteString(fmt.Sprintf("Milestone #%d (%s)\n\n", m.Number, strings.Title(m.State)))

	// Description
	if m.Description != "" {
		sb.WriteString(fmt.Sprintf("Description:\n  %s\n\n", m.Description))
	}

	// Due date
	if !m.DueOn.IsZero() {
		dueDate := m.DueOn.Format("January 2, 2006")
		daysUntil := int(time.Until(m.DueOn).Hours() / 24)

		if daysUntil > 0 {
			sb.WriteString(fmt.Sprintf("Due date: %s (%d days from now)\n\n", dueDate, daysUntil))
		} else if daysUntil == 0 {
			sb.WriteString(fmt.Sprintf("Due date: %s (today)\n\n", dueDate))
		} else {
			sb.WriteString(fmt.Sprintf("Due date: %s (%d days overdue)\n\n", dueDate, -daysUntil))
		}
	}

	// Progress
	total := m.OpenIssues + m.ClosedIssues
	if total > 0 {
		percentage := (float64(m.ClosedIssues) / float64(total)) * 100
		sb.WriteString(fmt.Sprintf("Progress:\n"))
		sb.WriteString(fmt.Sprintf("  %d of %d issues completed (%.0f%%)\n", m.ClosedIssues, total, percentage))
		sb.WriteString(fmt.Sprintf("  %s\n\n", progressBar(m.ClosedIssues, total)))
	} else {
		sb.WriteString("Progress:\n  No issues assigned\n\n")
	}

	// URL
	sb.WriteString(fmt.Sprintf("URL: %s\n", m.URL))

	return sb.String()
}

// formatDueDate formats a due date for table display
func formatDueDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}

	daysUntil := int(time.Until(t).Hours() / 24)

	if daysUntil > 0 {
		return t.Format("2006-01-02")
	} else if daysUntil == 0 {
		return "Today"
	} else {
		return fmt.Sprintf("%s ⚠", t.Format("2006-01-02"))
	}
}

// formatProgress formats progress as "closed/total (percentage%)"
func formatProgress(open, closed int) string {
	total := open + closed
	if total == 0 {
		return "0/0 (0%)"
	}
	percentage := (float64(closed) / float64(total)) * 100
	return fmt.Sprintf("%d/%d (%.0f%%)", closed, total, percentage)
}

// progressBar creates a visual progress bar
func progressBar(completed, total int) string {
	if total == 0 {
		return ""
	}

	width := 30
	filled := int((float64(completed) / float64(total)) * float64(width))

	bar := strings.Repeat("▓", filled) + strings.Repeat("░", width-filled)
	return bar
}

// truncate truncates a string to a maximum length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
