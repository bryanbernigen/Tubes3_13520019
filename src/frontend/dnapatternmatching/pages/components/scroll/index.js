import React, { useEffect, useState } from "react";
import styles from './Scroll.module.css';

const ScrollTop = () => {
  const [isScroll, setIsScroll] = useState(false);

  const scrollTop = () => {
    window.scrollTo({
      top: 0,
      behavior: "smooth"
    });
  };

  useEffect(() => {
    window.addEventListener("scroll", () => {
        if (window.pageYOffset > 100) {
            setIsScroll(true);
        } else {
            setIsScroll(false);
        }
    });
  }, []);

  return (
    <div className={isScroll ? styles.scrollContainer : styles.scrollContainerHidden}>
        <button onClick={scrollTop} className={styles.scrollBtn}>
            <img src="/arrowup.png" className={styles.scrollArrow}/>
        </button>
    </div>
  )
}

export default ScrollTop