import React, { useEffect, useState } from 'react';
import styles from "./Navbar.module.css"
import Link from 'next/link'

const Navbar = () => {
    const [navbar, setNavbar] = useState(false);

    useEffect(() => {
        const changeBackground = () => {
            if (window.scrollY >= 10) {
                setNavbar(true);
            }
            else {
                setNavbar(false);
            }
        }
        window.addEventListener('scroll', changeBackground);

        return () => window.removeEventListener("scroll", changeBackground);
    }, [])

    return (
        <nav className={navbar ? styles.navbarActive : styles.navbar} >
            <div className={styles.navContainer}>
                <a href='/'>
                    <img src="1.png" className={styles.logoNavbar} />
                </a>
                <div className={styles.menuContainer}>
                    <Link href="#HomeSection">
                        <a className={styles.menuItem}>
                            <span className={styles.menuItemText}>Home</span>
                        </a>
                    </Link>
                    <Link href="#AddDiseaseSection">
                        <a className={styles.menuItem}>
                            <span className={styles.menuItemText}>AddDisease</span>
                        </a>
                    </Link>
                    <Link href="#TestDNASection">
                        <a className={styles.menuItem}>
                            <span className={styles.menuItemText}>DNATest</span>
                        </a>
                    </Link>
                    <Link href="#SearchBarSection">
                        <a className={styles.menuItem}>
                            <span className={styles.menuItemText}>SearchResults</span>
                        </a>
                    </Link>
                    <Link href="#AboutUsSection">
                        <a className={styles.menuItem}>
                            <span className={styles.menuItemText}>About</span>
                        </a>
                    </Link>
                    <Link href="#ContactSection">
                        <a className={styles.menuItem}>
                            <span className={styles.menuItemText}>Contact</span>
                        </a>
                    </Link>
                    
                </div>
            </div>
        </nav>
    )
};

export default Navbar;