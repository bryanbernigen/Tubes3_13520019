import React from 'react';
import styles from './PredictionCard.module.css';

const PredictionCard = ({object}) => {
    return (
        <div className={styles.cardCt}>
            <span>{object.tanggalprediksi}</span>
            <span> — </span>
            <span>{object.namapasien}</span>
            <span> — </span>
            <span>{object.namapenyakit}</span>
            <span> — </span>
            <span>{object.statuspenyakit.toString()}</span>
        </div>
    )
}

export default PredictionCard;