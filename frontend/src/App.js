import { BrowserRouter, Routes, Route } from "react-router-dom";
import SearchPage from "./pages/search/SearchPage"
import Case from "./pages/case/Case"

export default function App() {
	return (
		<>
			<BrowserRouter>
				<Routes>
					<Route path="/" element={<SearchPage />} />
					<Route path="/case" element={<Case />} />
				</Routes>
			</BrowserRouter>
		</>
	);
}
