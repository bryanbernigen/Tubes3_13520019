import React from 'react';
import styles from './Footer.module.css';

const Footer = () => {
    return (
        <div className={styles.footerCt}>
            <span className={styles.footerText}>
                Contact us on
            </span>
            <div className={styles.contacts}>
                <div className={styles.githubContact}>
                    <a href='/'>
                        <button className={styles.githubButton}>
                            <img src='footer/GitHub-Mark-64px.png' className={styles.logo}/>
                            <span className={styles.name}>runnnayy</span>
                        </button>
                    </a>
                </div>
                <div className={styles.githubContact}>
                    <a href='/'>
                        <button className={styles.githubButton}>
                            <img src='footer/GitHub-Mark-64px.png' className={styles.logo}/>
                            <span className={styles.name}>bryanbernigen</span>
                        </button>
                    </a>
                </div>
                <div className={styles.githubContact}>
                    <a href='/'>
                        <button className={styles.githubButton}>
                            <img src='footer/GitHub-Mark-64px.png' className={styles.logo}/>
                            <span className={styles.name}>nk-kyle</span>
                        </button>
                    </a>
                </div>
            </div>
            <div className={styles.year}>
                <span className={styles.yearText}>
                    © 2022
                </span>
            </div>
        </div>
    )
};

export default Footer;