import { useLocation } from "react-router-dom";
import SearchBar from "../../components/SearchBar"
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
			<div className="search-bar-wrapper-case-page">
				<SearchBar />
			</div>
			<div style={{ margin: '16px 34px' }}>
				<h2 style={{ fontWeight: '400' }}>
					Processo n. { courtCase.court_case.cnj } do { courtCase.court_case.court_of_origin }
				</h2>
				<p style={{ margin: '4px 0' }}>
					Distribuído em { formatDate(courtCase.court_case.start_date) } | { courtCase.court_case.plaintiff } vs { courtCase.court_case.defendant }
				</p>
				<div className="updates-title-box">Movimentações</div>
				{ courtCase.court_case.updates.map((update) => (
					<div className="update-box-wrapper">
					<p style={{ marginBottom: '8px'}}>{ formatDate(update.update_date) }</p>
						<p>{ update.update_details }</p>
					</div>
				))}
			</div>
		</>
	)
}
