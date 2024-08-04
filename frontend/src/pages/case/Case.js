import { useLocation } from "react-router-dom";
import "./styles.css"

export default function Case() {
	const { state } = useLocation()
	const { courtCase } = state || {}

	function formatDate(dateStr) {
		const date = new Date(dateStr);

		const day = date.getUTCDate();
		const month = date.getUTCMonth() + 1;
		const year = date.getUTCFullYear();

		return `${day}/${month}/${year}`;
	}

	return (
		<>
			<h1>Processo n. { courtCase.court_case.cnj } do { courtCase.court_case.court_of_origin }</h1>
			<p>Distribuído em { formatDate(courtCase.court_case.start_date) }</p>
			<p>{ courtCase.court_case.plaintiff } vs { courtCase.court_case.defendant }</p>
			<div>Movimentações</div>
			{ courtCase.court_case.updates.map((update) => (
				<>
					<p>{ formatDate(update.update_date) }</p>
					<p>{ update.update_details }</p>
				</>
			))}
		</>
	)
}
