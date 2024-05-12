package email

const (
	crontabAlertHTML = `<!DOCTYPE html>
	<html lang="en">
	
	<head>
		<meta charset="UTF-8">
		<title>template</title>
		<style>
			table {
				border-collapse: collapse;
				border-spacing: 0;
			}
			thead {
				font-weight: bold;
			}
			tbody tr:nth-child(odd) {
				background: gray;
			}
			td {
				padding: 5px;
			}
		</style>
	</head>
	<body>
		<table border="1">
			<tbody>
				<tr>
					<td>taskName</td>
					<td>%s</td>
				</tr>
				<tr>
					<td>resultCode</td>
					<td>%d</td>
				</tr>
				<tr>
					<td>result</td>
					<td>%s</td>
				</tr>

				<tr>
					<td>viewURL</td>
					<td>%s</td>
				</tr>
			</tbody>
		</table>
	</body>
	
	</html>`
)
