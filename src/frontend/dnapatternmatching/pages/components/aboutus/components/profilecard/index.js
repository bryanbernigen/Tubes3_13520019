import React from 'react';
import styles from './ProfileCard.module.css';
import Typewriter from 'typewriter-effect';

const ProfileCard = ({ImgPath, ImgAlt, Name, NIM}) => {
    return (
        <div className={styles.profileCardBox}>
            <img src={ImgPath} alt={ImgAlt} className={styles.photo} />
            <span className={styles.name}>
                <Typewriter
                    onInit={(typewriter) => {
                        typewriter.typeString(Name)
                            .pauseFor(2500)
                            .start()
                            .deleteAll();
                    }}
                    options={{
                        autoStart: true,
                        loop: true,
                    }}
                />
            </span>
            <span className={styles.nim}>
                <Typewriter
                    onInit={(typewriter) => {
                        typewriter.typeString(NIM)
                            .pauseFor(2500)
                            .start()
                            .deleteAll();
                    }}
                    options={{
                        autoStart: true,
                        loop: true,
                    }}
                />
            </span>
        </div>
    )
}

export default ProfileCard;