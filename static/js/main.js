$(function() {
    $.getJSON("/tasks", function(tasks) {
        var $tbody = $("<tbody></tbody>");
        Object.keys(tasks).forEach(function(taskName) {
            var $tr = $tbody.append("<tr></tr>");
            var task = tasks[taskName];
            $tr.append("<td>" + task.name + "</td>");
            $tr.append("<td>" + task.children + "</td>");
            $tr.append("<td>" + task.parents + "</td>");
            var buttonClasses = "btn btn-default triggerTask";
            if (task.parents.length > 0) {
                buttonClasses += " disabled"
            }
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
        })
    });

    $("#registerNewTask").click(function() {
        var taskName = $("#taskName").val();
        var children = $("#taskChildren").val().split(/,\s?/).filter(function(taskName) { return taskName !== "" });
        var parents = $("#taskParents").val().split(/,\s?/).filter(function(taskName) { return taskName !== "" });
        $.ajax({
            type: "POST",
            url: "/tasks/register",
            data: JSON.stringify({
                "name": taskName,
                "children": children,
                "parents": parents
            }),
            dataType: "json"
        });
    })
});
