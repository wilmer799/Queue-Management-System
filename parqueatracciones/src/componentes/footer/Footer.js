import React from 'react'

class Footer extends React.Component {
    render() {
        return (
            <footer className="fixed-bottom bg-primary" >
                <p className="terminos">&copy; {(new Date().getFullYear())} Go Aventura, Inc. &middot; <a  className="text-white" href="#">Política de Privacidad</a> &middot; <a  className="text-white" href="#">Términos</a></p>
            </footer>

        )
    }
}

export default Footer;

