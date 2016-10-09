$(function() {
    $.getJSON("/tasks", function(tasks) {
        var $tbody = $("<tbody></tbody>");
        Object.keys(tasks).forEach(function(taskName) {
            var $tr = $tbody.append("<tr></tr>");
            var task = tasks[taskName];
            var env = task.env === null ? "none" : Object.keys(task.env).map(function(key) { return key + "=" + task.env[key] }).join(",");
            $tr.append("<td>" + task.name + "</td>");
            $tr.append("<td>" + env + "</td>");
            $tr.append("<td>" + task.children + "</td>");
            $tr.append("<td>" + task.parents + "</td>");
            var buttonClasses = "btn btn-default triggerTask";
            if (task.parents.length > 0) {
                buttonClasses += " disabled"
            }
            $tr.append('<button data-task-name="' + task.name + '" ' +
                'data-task-env="' + env + '" ' +
                'data-task-children="' + task.children + '" ' +
                'data-task-parents="' + task.parents + '" ' +
                'class="btn btn-default editTask" type="submit">edit</button>');
            $tr.append('<button data-task-name="' + task.name + '" class="' + buttonClasses + '" type="submit">trigger</button>');

        });

        $("#registeredTasks").append($tbody);
        $(".triggerTask").click(function() {
            var taskName = $(this).attr('data-task-name');
            $.ajax({
                type: "POST",
                url: "/tasks/trigger",
                data: JSON.stringify({
                    "task": taskName
                }),
                dataType: "json"
            });
        });

        $(".editTask").click(function() {
            var name = $(this).attr('data-task-name');
            var env = $(this).attr('data-task-env');
            var children = $(this).attr('data-task-children');
            var parents = $(this).attr('data-task-parents');

            $("#taskName").val(name);
            $("#taskEnv").val(env);
            $("#taskChildren").val(children);
            $("#taskParents").val(parents);

            $("#registerNewTaskModal").modal('show');
        })
    });

    $("#registerNewTask").click(function() {
        var taskName = $("#taskName").val();
        var children = $("#taskChildren").val().split(/,\s?/).filter(function(taskName) { return taskName !== "" });
        var parents = $("#taskParents").val().split(/,\s?/).filter(function(taskName) { return taskName !== "" });
        var env = $("#taskEnv").val().split(",").reduce(function(acc, keyValue) {
            var keyValueArr = keyValue.split("=");
            if (keyValueArr.length !== 2) {
                return acc
            }
            acc[keyValueArr[0]] = keyValueArr[1];
            debugger;
            return acc
        }, {});

        $.ajax({
            type: "POST",
            url: "/tasks/register",
            data: JSON.stringify({
                "name": taskName,
                "children": children,
                "parents": parents,
                "env": env
            }),
            dataType: "json",
            success: function() {
              location.reload()
            }
        });
    })
});
