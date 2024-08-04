import { BrowserRouter, Routes, Route } from "react-router-dom";
import SearchPage from "./pages/search/SearchPage"
import Case from "./pages/case/Case"
import RegisterCase from "./pages/register_case/RegisterCase"

export default function App() {
	return (
		<>
			<BrowserRouter>
				<Routes>
					<Route path="/" element={<SearchPage />} />
					<Route path="/case" element={<Case />} />
					<Route path="/register_case" element={<RegisterCase />} />
				</Routes>
			</BrowserRouter>
		</>
	);
}
