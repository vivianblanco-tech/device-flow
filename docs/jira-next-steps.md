Next Steps to Integrate

Add environment variables to .env:
   JIRA_URL=https://bairesdev.atlassian.net
   JIRA_CLIENT_ID=your-client-id
   JIRA_CLIENT_SECRET=your-client-secret

Initialize JIRA client in cmd/web/main.go

Update your handlers to call JIRA methods:
- In pickup form: Create ticket
- In status updates: Sync status
- In delivery form: Add comment

Add database field (optional):
   ALTER TABLE shipments ADD COLUMN jira_ticket_key VARCHAR(50);