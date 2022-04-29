import React, { useEffect } from 'react';
import styles from './Home.module.css';
import Typewriter from 'typewriter-effect';

const HomePage = () => {
    return (
        <div className={styles.homeContainer}>
            <div className={styles.homeBG}>
                <div className={styles.homeTitle}>
                    <Typewriter
                        onInit={(typewriter) => {
                            typewriter.typeString('DNA Pattern Matching')
                                .pauseFor(2500)
                                .start()
                                .deleteAll();
                        }}
                        options={{
                            autoStart: true,
                            loop: true,
                        }}
                    />
                </div>
            </div>
        </div>
    )
}

export default HomePage