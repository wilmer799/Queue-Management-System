import React from 'react'
class Menu extends React.Component {
    render() {
        return (
			<div>
				<header>
					<nav className="navbar navbar-expand-md navbar-dark fixed-top bg-primary">
						<a className="navbar-brand" href="#">Go Aventura</a>
						<button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
							<span className="navbar-toggler-icon"></span>
						</button>
						<div className="collapse navbar-collapse" id="navbarCollapse">
							<ul className="navbar-nav mr-auto">
								<li className="nav-item active">
									<a className="nav-link" href="#">Mapa </a>
								</li>
								<li className="nav-item">
									<a className="nav-link" href="#">Estado Visitantes</a>
								</li>
							</ul>
						</div>
					</nav>
				</header>
			</div>
        )
    }
}

export default Menu;