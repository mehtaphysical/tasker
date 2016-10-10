$(function() {
    fetchHistory();
});

function fetchHistory() {
    $.getJSON("/tasks/history", function(tasks) {
        var $tbody = $('<tbody></tbody>');
        tasks.forEach(function(task) {
            var rowId = task.id + "-" + task.name.replace(/\s+/g, "_");
            var $tr = $('<tr id="' + rowId + '"></tr>');
            var $outputButton = $('<td><button id="output-' + rowId + '" class="btn btn-primary">Output</button></td>');

            $tr.append('<td>' + task.id + '</td>');
            $tr.append('<td>' + task.name + '</td>');
            $tr.append('<td>' + task.status + '</td>');
            $tr.append($outputButton);
            $tbody.append($tr);

            $outputButton.click(function() {
                $("#taskOutput").text(task.output);
                $("#taskHistoryModal").modal('show');
            })
        });
        $('#taskHistory').append($tbody);
    });
}
