# chmgt
Easy change management software for the masses.

## Ideas for possible features
-	A method of assigning or locking task items
    -   gives a clear indicator of who should work on an item
    -   prevents changing an item once work has begun
-	Reminder notifications of items that are due soon
-	Summary of items scheduled
    -   Customizable schedule for sending email
    -   Customizable schedule for what items should be included in email (how far to look ahead)
-	Blackout times for maintenance on services
    -   Disallow scheduling of work items without override
    -   Notification of items already scheduled when entering blackouts
-	Dependency tracking for changes
    -   What other changes are required prior to current change
-	Breakout of steps and step completion indication
    -   Allow a single change request to be broken into multiple tasks
    -   Each task can be assigned to a different person
    -   Indicator for the number of tasks completed out of total (possible percentage as well)
-   Plugins to allow the change management tool make the changes for you
    -   GitLab API to view and accept changes
    -   Ansible AWX API to activate the changes
    -   A scheduler tool to... ummm... schedule the changes
    -   If the change doesn't end with a sucessful return code,          automatically revert the change and send a notification with the logs
