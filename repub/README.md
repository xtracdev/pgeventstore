Register task:

<pre>
aws ecs register-task-definition --cli-input-json file://$PWD/taskdef.json
</pre>

Run task:

<pre>
aws ecs run-task --cluster green --task-definition repub
</pre>