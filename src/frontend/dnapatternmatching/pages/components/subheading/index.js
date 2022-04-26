import React from 'react';
import styles from './Subheading.module.css';

const Subheading = ({Text, Color}) => {
    return (
        <div className={styles.subheadingContainer}>
            <span className={styles.subheadingText}>
                {
                    Color === "black" 
                        ? <span className={styles.subheadingTextBlack}>{Text}</span>
                        : <span className={styles.subheadingTextWhite}>{Text}</span>
                }
            </span>
            <hr className={styles.thickLine} />
        </div>
    )
};

export default Subheading;