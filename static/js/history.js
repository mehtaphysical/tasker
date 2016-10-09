$(function() {
    fetchHistory();
});

function fetchHistory() {
    $.getJSON("/tasks/history", function(tasks) {
        var $tbody = $('<tbody></tbody>');
        tasks.forEach(function(task) {
            var $tr = $('<tr></tr>');
            $tr.append('<td>' + task.id + '</td>');
            $tr.append('<td>' + task.name + '</td>');
            $tr.append('<td>' + task.status + '</td>');
            $tbody.append($tr);
        });
        $('#taskHistory').append($tbody);
    });
}
