import { useNavigate } from "react-router-dom";
import SearchBar from "../../components/SearchBar"
import "./styles.css"

export default function SearchPage() {
	const navigate = useNavigate()

	return (
		<div className="search-page-wrapper">
			<h1 className="search-page-title">Buscar</h1>
			<h3 style={{ textAlign: 'center', fontWeight: '400' }}>Busque um processo a partir do número unificado</h3>
			<SearchBar />
			<div className="div-center">
			<button className="pretty-btn" style={{ marginTop: '8px'}} onClick={() => navigate("/register_case")}>Register Case</button>
			</div>
		</div>
	)
}
