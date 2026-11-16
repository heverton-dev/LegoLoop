# Relatório Consolidado Web-QA — legoloop-studio-v2

- **Data:** 2026-11-15
- **URL Alvo:** http://127.0.0.1:7893
- **Status Final:** FAILED
- **Rotas Auditadas:** /, /inventario, /swarm

## Achados Consolidados

1. **[playwright]** (critical) `/ (desktop)`: Failed to load resource: net::ERR_CONNECTION_REFUSED
2. **[axe-core]** (critical) `/ (desktop)`: aria-required-children: Certain ARIA roles must contain particular children (1 elemento(s))
3. **[playwright]** (critical) `/inventario (desktop)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7893/inventario
Call log:
[2m  - navigating to "http://127.0.0.1:7893/inventario", waiting until "networkidle"[22m

4. **[playwright]** (critical) `/swarm (desktop)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7893/swarm
Call log:
[2m  - navigating to "http://127.0.0.1:7893/swarm", waiting until "networkidle"[22m

5. **[playwright]** (critical) `/ (mobile)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7893/
Call log:
[2m  - navigating to "http://127.0.0.1:7893/", waiting until "networkidle"[22m

6. **[playwright]** (critical) `/inventario (mobile)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7893/inventario
Call log:
[2m  - navigating to "http://127.0.0.1:7893/inventario", waiting until "networkidle"[22m

7. **[playwright]** (critical) `/swarm (mobile)`: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7893/swarm
Call log:
[2m  - navigating to "http://127.0.0.1:7893/swarm", waiting until "networkidle"[22m

8. **[renderlint]** (critical) `http://127.0.0.1:7893/`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  / @768px light: 88 finding(s) in 6768ms
  / @1280px light: 38 finding(s) in 6798ms
  / @375px dark: 84 finding(s) in 5510ms
  / @375px light: 92 finding(s) in 10084ms
  / @1280px dark: 38 finding(s) in 5506ms
  / @768px dark: 61 finding(s) in 5610ms
  /swarm @375px light: 5 finding(s) in 1146ms
  /inventario @375px light: 3 finding(s) in 6298ms
  /swarm @1280px light: 0 finding(s) in 75ms
  /inventario @768px light: 3 finding(s) in 4729ms
  /inventario @375px dark: 0 finding(s) in 61ms
  /inventario @768px dark: 0 finding(s) in 54ms
  /inventario @1280px dark: 0 finding(s) in 53ms
  /swarm @375px dark: 0 finding(s) in 66ms
  /swarm @768px dark: 0 finding(s) in 60ms
  /swarm @1280px dark: 0 finding(s) in 63ms
  /swarm @768px light: 0 finding(s) in 1943ms
  data stress skipped on http://127.0.0.1:7893/inventario: page.reload: net::ERR_CONNECTION_REFUSED
  /inventario @1280px light: 3 finding(s) in 3474ms
9. **[renderlint]** (critical) `http://127.0.0.1:7893/inventario`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /inventario @768px light: 3 finding(s) in 6135ms
  /inventario @1280px light: 3 finding(s) in 6284ms
  /inventario @375px dark: 3 finding(s) in 8392ms
  /inventario @768px dark: 3 finding(s) in 5167ms
  /inventario @1280px dark: 3 finding(s) in 5110ms
  data stress skipped on http://127.0.0.1:7893/: page.reload: net::ERR_CONNECTION_REFUSED
  data stress skipped on http://127.0.0.1:7893/: page.reload: net::ERR_CONNECTION_REFUSED
  / @1280px light: 50 finding(s) in 2973ms
  loading-state check skipped on http://127.0.0.1:7893/: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7893/
  / @375px light: 93 finding(s) in 3212ms
  /swarm @375px light: 0 finding(s) in 91ms
  /swarm @768px light: 0 finding(s) in 67ms
  /swarm @1280px light: 0 finding(s) in 75ms
  / @375px dark: 0 finding(s) in 97ms
  / @768px dark: 0 finding(s) in 82ms
  /inventario @375px light: 4 finding(s) in 14658ms
  / @1280px dark: 0 finding(s) in 66ms
  data stress skipped on http://127.0.0.1:7893/: page.reload: net::ERR_CONNECTION_REFUSED
  / @768px light: 89 finding(s) in 3390ms
10. **[renderlint]** (critical) `http://127.0.0.1:7893/swarm`: Checking 1 page(s) (+ crawling links) at 375/768/1280px…
  /swarm @1280px light: 9 finding(s) in 9979ms
  /swarm @768px light: 9 finding(s) in 9994ms
  /swarm @375px dark: 8 finding(s) in 5202ms
  /swarm @375px light: 9 finding(s) in 11854ms
  data stress skipped on http://127.0.0.1:7893/: page.reload: net::ERR_CONNECTION_REFUSED
  loading-state check skipped on http://127.0.0.1:7893/: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7893/
  / @375px light: 93 finding(s) in 2555ms
  data stress skipped on http://127.0.0.1:7893/: page.reload: net::ERR_CONNECTION_REFUSED
  / @768px light: 92 finding(s) in 2396ms
  / @1280px light: 0 finding(s) in 78ms
  /inventario @375px light: 0 finding(s) in 89ms
  /inventario @768px light: 0 finding(s) in 80ms
  /inventario @1280px light: 0 finding(s) in 91ms
  / @375px dark: 0 finding(s) in 70ms
  / @768px dark: 0 finding(s) in 66ms
  / @1280px dark: 0 finding(s) in 69ms
  /swarm @768px dark: 1 finding(s) in 4588ms
  /swarm @1280px dark: 9 finding(s) in 5104ms
11. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f4f3ef" (linha 18); use var(--*) de tokens.css
12. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ffffff" (linha 19); use var(--*) de tokens.css
13. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#eeece6" (linha 20); use var(--*) de tokens.css
14. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#e3dfd6" (linha 21); use var(--*) de tokens.css
15. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#d6d1c6" (linha 22); use var(--*) de tokens.css
16. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#9f998c" (linha 23); use var(--*) de tokens.css
17. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1a1917" (linha 24); use var(--*) de tokens.css
18. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#524e46" (linha 25); use var(--*) de tokens.css
19. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#d9d5ca" (linha 26); use var(--*) de tokens.css
20. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b83908" (linha 27); use var(--*) de tokens.css
21. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ffffff" (linha 28); use var(--*) de tokens.css
22. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fde8de" (linha 29); use var(--*) de tokens.css
23. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1f6f31" (linha 30); use var(--*) de tokens.css
24. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ffffff" (linha 31); use var(--*) de tokens.css
25. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#e6f6ea" (linha 32); use var(--*) de tokens.css
26. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b91c1c" (linha 33); use var(--*) de tokens.css
27. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fde8e8" (linha 34); use var(--*) de tokens.css
28. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b91c1c" (linha 35); use var(--*) de tokens.css
29. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ffffff" (linha 36); use var(--*) de tokens.css
30. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#52585f" (linha 37); use var(--*) de tokens.css
31. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#eceef0" (linha 38); use var(--*) de tokens.css
32. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6f6a5e" (linha 39); use var(--*) de tokens.css
33. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.06)" (linha 40); use var(--*) de tokens.css
34. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.06)" (linha 40); use var(--*) de tokens.css
35. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#121210" (linha 55); use var(--*) de tokens.css
36. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1b1b19" (linha 56); use var(--*) de tokens.css
37. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#242420" (linha 57); use var(--*) de tokens.css
38. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2e2e29" (linha 58); use var(--*) de tokens.css
39. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#363530" (linha 59); use var(--*) de tokens.css
40. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#57544c" (linha 60); use var(--*) de tokens.css
41. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f0efe9" (linha 61); use var(--*) de tokens.css
42. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b0aca0" (linha 62); use var(--*) de tokens.css
43. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2a2926" (linha 63); use var(--*) de tokens.css
44. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ff8247" (linha 64); use var(--*) de tokens.css
45. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1a0b04" (linha 65); use var(--*) de tokens.css
46. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3a2116" (linha 66); use var(--*) de tokens.css
47. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#4ade80" (linha 67); use var(--*) de tokens.css
48. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#042f1e" (linha 68); use var(--*) de tokens.css
49. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#142c1d" (linha 69); use var(--*) de tokens.css
50. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ff7b7b" (linha 70); use var(--*) de tokens.css
51. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3b1919" (linha 71); use var(--*) de tokens.css
52. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#c92a2a" (linha 72); use var(--*) de tokens.css
53. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ffffff" (linha 73); use var(--*) de tokens.css
54. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#a8aeb4" (linha 74); use var(--*) de tokens.css
55. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#26282a" (linha 75); use var(--*) de tokens.css
56. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#878276" (linha 76); use var(--*) de tokens.css
57. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.4)" (linha 77); use var(--*) de tokens.css
58. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.35)" (linha 77); use var(--*) de tokens.css
59. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#121210" (linha 83); use var(--*) de tokens.css
60. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1b1b19" (linha 84); use var(--*) de tokens.css
61. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#242420" (linha 85); use var(--*) de tokens.css
62. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2e2e29" (linha 86); use var(--*) de tokens.css
63. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#363530" (linha 87); use var(--*) de tokens.css
64. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#57544c" (linha 88); use var(--*) de tokens.css
65. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f0efe9" (linha 89); use var(--*) de tokens.css
66. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b0aca0" (linha 90); use var(--*) de tokens.css
67. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2a2926" (linha 91); use var(--*) de tokens.css
68. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ff8247" (linha 92); use var(--*) de tokens.css
69. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1a0b04" (linha 93); use var(--*) de tokens.css
70. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3a2116" (linha 94); use var(--*) de tokens.css
71. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#4ade80" (linha 95); use var(--*) de tokens.css
72. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#042f1e" (linha 96); use var(--*) de tokens.css
73. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#142c1d" (linha 97); use var(--*) de tokens.css
74. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ff7b7b" (linha 98); use var(--*) de tokens.css
75. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3b1919" (linha 99); use var(--*) de tokens.css
76. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#c92a2a" (linha 100); use var(--*) de tokens.css
77. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ffffff" (linha 101); use var(--*) de tokens.css
78. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#a8aeb4" (linha 102); use var(--*) de tokens.css
79. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#26282a" (linha 103); use var(--*) de tokens.css
80. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#878276" (linha 104); use var(--*) de tokens.css
81. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.4)" (linha 105); use var(--*) de tokens.css
82. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.35)" (linha 105); use var(--*) de tokens.css
83. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.22)" (linha 194); use var(--*) de tokens.css
84. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(14, 165, 233, 0.12)" (linha 280); use var(--*) de tokens.css
85. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0284c7" (linha 280); use var(--*) de tokens.css
86. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0369a1" (linha 280); use var(--*) de tokens.css
87. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(14, 165, 233, 0.22)" (linha 281); use var(--*) de tokens.css
88. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0284c7" (linha 281); use var(--*) de tokens.css
89. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#38bdf8" (linha 282); use var(--*) de tokens.css
90. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(56, 189, 248, 0.55)" (linha 282); use var(--*) de tokens.css
91. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#38bdf8" (linha 284); use var(--*) de tokens.css
92. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(56, 189, 248, 0.55)" (linha 284); use var(--*) de tokens.css
93. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0284c7" (linha 600); use var(--*) de tokens.css
94. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0284c7" (linha 651); use var(--*) de tokens.css
95. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0284c7" (linha 673); use var(--*) de tokens.css
96. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(14, 165, 233, 0.12)" (linha 674); use var(--*) de tokens.css
97. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f59e0b" (linha 707); use var(--*) de tokens.css
98. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(245, 158, 11, 0.22)" (linha 708); use var(--*) de tokens.css
99. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fbbf24" (linha 711); use var(--*) de tokens.css
100. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(245, 158, 11, 0.35)" (linha 712); use var(--*) de tokens.css
101. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fbbf24" (linha 712); use var(--*) de tokens.css
102. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(245, 158, 11, 0.08)" (linha 715); use var(--*) de tokens.css
103. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(245, 158, 11, 0.2)" (linha 716); use var(--*) de tokens.css
104. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#9a3412" (linha 719); use var(--*) de tokens.css
105. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fbbf24" (linha 722); use var(--*) de tokens.css
106. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fbbf24" (linha 726); use var(--*) de tokens.css
107. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f59e0b" (linha 730); use var(--*) de tokens.css
108. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(14, 165, 233, 0.14)" (linha 783); use var(--*) de tokens.css
109. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0369a1" (linha 783); use var(--*) de tokens.css
110. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0284c7" (linha 783); use var(--*) de tokens.css
111. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#38bdf8" (linha 784); use var(--*) de tokens.css
112. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#38bdf8" (linha 786); use var(--*) de tokens.css
113. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.45)" (linha 903); use var(--*) de tokens.css
114. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.65)" (linha 1184); use var(--*) de tokens.css
115. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(255 255 255 / 0.05)" (linha 1184); use var(--*) de tokens.css
116. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(0 0 0 / 0.7)" (linha 1200); use var(--*) de tokens.css
117. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(239 68 68 / 0.08)" (linha 1315); use var(--*) de tokens.css
118. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ef4444" (linha 1316); use var(--*) de tokens.css
119. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ef4444" (linha 1338); use var(--*) de tokens.css
120. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fff" (linha 1422); use var(--*) de tokens.css
121. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgb(59 130 246 / 0.15)" (linha 1763); use var(--*) de tokens.css
122. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(0, 0, 0, 0.15)" (linha 1901); use var(--*) de tokens.css
123. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fff" (linha 1905); use var(--*) de tokens.css
124. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f59e0b" (linha 2269); use var(--*) de tokens.css
125. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(245, 158, 11, 0.1)" (linha 2270); use var(--*) de tokens.css
126. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#10b981" (linha 2273); use var(--*) de tokens.css
127. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(16, 185, 129, 0.1)" (linha 2274); use var(--*) de tokens.css
128. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ef4444" (linha 2277); use var(--*) de tokens.css
129. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(239, 68, 68, 0.1)" (linha 2278); use var(--*) de tokens.css
130. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0f172a" (linha 2281); use var(--*) de tokens.css
131. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#38bdf8" (linha 2282); use var(--*) de tokens.css
132. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(0, 0, 0, 0.2)" (linha 2322); use var(--*) de tokens.css
133. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f97316" (linha 2376); use var(--*) de tokens.css
134. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f97316" (linha 2376); use var(--*) de tokens.css
135. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(249, 115, 22, 0.18)" (linha 2376); use var(--*) de tokens.css
136. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#eab308" (linha 2377); use var(--*) de tokens.css
137. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#eab308" (linha 2377); use var(--*) de tokens.css
138. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(234, 179, 8, 0.16)" (linha 2377); use var(--*) de tokens.css
139. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#8b5cf6" (linha 2378); use var(--*) de tokens.css
140. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#8b5cf6" (linha 2378); use var(--*) de tokens.css
141. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(139, 92, 246, 0.16)" (linha 2378); use var(--*) de tokens.css
142. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#06b6d4" (linha 2379); use var(--*) de tokens.css
143. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#06b6d4" (linha 2379); use var(--*) de tokens.css
144. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(6, 182, 212, 0.16)" (linha 2379); use var(--*) de tokens.css
145. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2563eb" (linha 2380); use var(--*) de tokens.css
146. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2563eb" (linha 2380); use var(--*) de tokens.css
147. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(37, 99, 235, 0.16)" (linha 2380); use var(--*) de tokens.css
148. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#10b981" (linha 2381); use var(--*) de tokens.css
149. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#10b981" (linha 2381); use var(--*) de tokens.css
150. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(16, 185, 129, 0.16)" (linha 2381); use var(--*) de tokens.css
151. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ec4899" (linha 2382); use var(--*) de tokens.css
152. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ec4899" (linha 2382); use var(--*) de tokens.css
153. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(236, 72, 153, 0.16)" (linha 2382); use var(--*) de tokens.css
154. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ef4444" (linha 2383); use var(--*) de tokens.css
155. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ef4444" (linha 2383); use var(--*) de tokens.css
156. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(239, 68, 68, 0.16)" (linha 2383); use var(--*) de tokens.css
157. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#9a3412" (linha 2385); use var(--*) de tokens.css
158. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(249, 115, 22, 0.18)" (linha 2385); use var(--*) de tokens.css
159. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(249, 115, 22, 0.35)" (linha 2385); use var(--*) de tokens.css
160. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#854d0e" (linha 2386); use var(--*) de tokens.css
161. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(234, 179, 8, 0.20)" (linha 2386); use var(--*) de tokens.css
162. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(234, 179, 8, 0.35)" (linha 2386); use var(--*) de tokens.css
163. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#5b21b6" (linha 2387); use var(--*) de tokens.css
164. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(139, 92, 246, 0.17)" (linha 2387); use var(--*) de tokens.css
165. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(139, 92, 246, 0.34)" (linha 2387); use var(--*) de tokens.css
166. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0e7490" (linha 2388); use var(--*) de tokens.css
167. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(6, 182, 212, 0.17)" (linha 2388); use var(--*) de tokens.css
168. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(6, 182, 212, 0.34)" (linha 2388); use var(--*) de tokens.css
169. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1d4ed8" (linha 2389); use var(--*) de tokens.css
170. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(37, 99, 235, 0.17)" (linha 2389); use var(--*) de tokens.css
171. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(37, 99, 235, 0.34)" (linha 2389); use var(--*) de tokens.css
172. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#047857" (linha 2390); use var(--*) de tokens.css
173. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(16, 185, 129, 0.17)" (linha 2390); use var(--*) de tokens.css
174. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(16, 185, 129, 0.34)" (linha 2390); use var(--*) de tokens.css
175. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#be185d" (linha 2391); use var(--*) de tokens.css
176. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(236, 72, 153, 0.17)" (linha 2391); use var(--*) de tokens.css
177. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(236, 72, 153, 0.34)" (linha 2391); use var(--*) de tokens.css
178. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b91c1c" (linha 2392); use var(--*) de tokens.css
179. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(239, 68, 68, 0.17)" (linha 2392); use var(--*) de tokens.css
180. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(239, 68, 68, 0.34)" (linha 2392); use var(--*) de tokens.css
181. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(249, 115, 22, 0.14)" (linha 2394); use var(--*) de tokens.css
182. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#9a3412" (linha 2394); use var(--*) de tokens.css
183. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(234, 179, 8, 0.16)" (linha 2395); use var(--*) de tokens.css
184. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#854d0e" (linha 2395); use var(--*) de tokens.css
185. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(139, 92, 246, 0.14)" (linha 2396); use var(--*) de tokens.css
186. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#5b21b6" (linha 2396); use var(--*) de tokens.css
187. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(6, 182, 212, 0.14)" (linha 2397); use var(--*) de tokens.css
188. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0e7490" (linha 2397); use var(--*) de tokens.css
189. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(37, 99, 235, 0.14)" (linha 2398); use var(--*) de tokens.css
190. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1d4ed8" (linha 2398); use var(--*) de tokens.css
191. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(16, 185, 129, 0.14)" (linha 2399); use var(--*) de tokens.css
192. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#047857" (linha 2399); use var(--*) de tokens.css
193. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(236, 72, 153, 0.14)" (linha 2400); use var(--*) de tokens.css
194. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#be185d" (linha 2400); use var(--*) de tokens.css
195. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(239, 68, 68, 0.14)" (linha 2401); use var(--*) de tokens.css
196. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b91c1c" (linha 2401); use var(--*) de tokens.css
197. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fdba74" (linha 2404); use var(--*) de tokens.css
198. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fde047" (linha 2406); use var(--*) de tokens.css
199. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#c4b5fd" (linha 2408); use var(--*) de tokens.css
200. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#67e8f9" (linha 2410); use var(--*) de tokens.css
201. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#93c5fd" (linha 2412); use var(--*) de tokens.css
202. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6ee7b7" (linha 2414); use var(--*) de tokens.css
203. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f9a8d4" (linha 2416); use var(--*) de tokens.css
204. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fca5a5" (linha 2418); use var(--*) de tokens.css
205. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fdba74" (linha 2421); use var(--*) de tokens.css
206. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fde047" (linha 2423); use var(--*) de tokens.css
207. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#c4b5fd" (linha 2425); use var(--*) de tokens.css
208. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#67e8f9" (linha 2427); use var(--*) de tokens.css
209. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#93c5fd" (linha 2429); use var(--*) de tokens.css
210. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6ee7b7" (linha 2431); use var(--*) de tokens.css
211. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f9a8d4" (linha 2433); use var(--*) de tokens.css
212. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fca5a5" (linha 2435); use var(--*) de tokens.css
213. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1d4ed8" (linha 2503); use var(--*) de tokens.css
214. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#60a5fa" (linha 2516); use var(--*) de tokens.css
215. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#60a5fa" (linha 2520); use var(--*) de tokens.css
216. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(0, 0, 0, 0.04)" (linha 2536); use var(--*) de tokens.css
217. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(59, 130, 246, 0.15)" (linha 2561); use var(--*) de tokens.css
218. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3b82f6" (linha 2562); use var(--*) de tokens.css
219. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(16, 185, 129, 0.15)" (linha 2566); use var(--*) de tokens.css
220. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#10b981" (linha 2567); use var(--*) de tokens.css
221. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(239, 68, 68, 0.15)" (linha 2571); use var(--*) de tokens.css
222. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ef4444" (linha 2572); use var(--*) de tokens.css
223. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(0, 0, 0, 0.03)" (linha 2580); use var(--*) de tokens.css
224. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(0, 0, 0, 0.05)" (linha 2614); use var(--*) de tokens.css
225. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(37, 99, 235, 0.22)" (linha 2632); use var(--*) de tokens.css
226. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(37, 99, 235, 0.10)" (linha 2632); use var(--*) de tokens.css
227. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3b82f6" (linha 2637); use var(--*) de tokens.css
228. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3b82f6" (linha 2638); use var(--*) de tokens.css
229. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(59, 130, 246, 0.35)" (linha 2638); use var(--*) de tokens.css
230. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(59, 130, 246, 0.18)" (linha 2638); use var(--*) de tokens.css
231. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f59e0b" (linha 2681); use var(--*) de tokens.css
232. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3b82f6" (linha 2685); use var(--*) de tokens.css
233. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ef4444" (linha 2689); use var(--*) de tokens.css
234. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#e5484d" (linha 2742); use var(--*) de tokens.css
235. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(255, 255, 255, 0.02)" (linha 2781); use var(--*) de tokens.css
236. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f1f5f9" (linha 2801); use var(--*) de tokens.css
237. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(255, 255, 255, 0.05)" (linha 2825); use var(--*) de tokens.css
238. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(255, 255, 255, 0.03)" (linha 2834); use var(--*) de tokens.css
239. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(255, 255, 255, 0.02)" (linha 2852); use var(--*) de tokens.css
240. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f1f5f9" (linha 2866); use var(--*) de tokens.css
241. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#e5484d" (linha 2869); use var(--*) de tokens.css
242. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(229, 72, 77, 0.08)" (linha 2870); use var(--*) de tokens.css
243. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#e5484d" (linha 2876); use var(--*) de tokens.css
244. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f1f5f9" (linha 2879); use var(--*) de tokens.css
245. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f1f5f9" (linha 2897); use var(--*) de tokens.css
246. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "rgba(255, 255, 255, 0.01)" (linha 2915); use var(--*) de tokens.css
247. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#38bdf8" (linha 2926); use var(--*) de tokens.css
248. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#38bdf8" (linha 2934); use var(--*) de tokens.css
249. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f1f5f9" (linha 2939); use var(--*) de tokens.css
250. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0b0f17" (linha 2956); use var(--*) de tokens.css
251. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f1f5f9" (linha 2962); use var(--*) de tokens.css
252. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#dbeafe" (linha 2981); use var(--*) de tokens.css
253. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1e40af" (linha 2982); use var(--*) de tokens.css
254. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#93c5fd" (linha 2983); use var(--*) de tokens.css
255. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#ede9fe" (linha 2987); use var(--*) de tokens.css
256. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#5b21b6" (linha 2988); use var(--*) de tokens.css
257. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#c4b5fd" (linha 2989); use var(--*) de tokens.css
258. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#d1fae5" (linha 2993); use var(--*) de tokens.css
259. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#065f46" (linha 2994); use var(--*) de tokens.css
260. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6ee7b7" (linha 2995); use var(--*) de tokens.css
261. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fef3c7" (linha 2999); use var(--*) de tokens.css
262. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#92400e" (linha 3000); use var(--*) de tokens.css
263. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fcd34d" (linha 3001); use var(--*) de tokens.css
264. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fce7f3" (linha 3005); use var(--*) de tokens.css
265. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#9d174d" (linha 3006); use var(--*) de tokens.css
266. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f9a8d4" (linha 3007); use var(--*) de tokens.css
267. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#cffafe" (linha 3011); use var(--*) de tokens.css
268. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#155e75" (linha 3012); use var(--*) de tokens.css
269. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#67e8f9" (linha 3013); use var(--*) de tokens.css
270. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#e0e7ff" (linha 3017); use var(--*) de tokens.css
271. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#3730a3" (linha 3018); use var(--*) de tokens.css
272. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#a5b4fc" (linha 3019); use var(--*) de tokens.css
273. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#e0e7ff" (linha 3024); use var(--*) de tokens.css
274. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#312e81" (linha 3025); use var(--*) de tokens.css
275. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#a5b4fc" (linha 3026); use var(--*) de tokens.css
276. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#172554" (linha 3031); use var(--*) de tokens.css
277. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#93c5fd" (linha 3032); use var(--*) de tokens.css
278. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1d4ed8" (linha 3033); use var(--*) de tokens.css
279. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2e1065" (linha 3037); use var(--*) de tokens.css
280. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#d8b4fe" (linha 3038); use var(--*) de tokens.css
281. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6d28d9" (linha 3039); use var(--*) de tokens.css
282. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#064e3b" (linha 3043); use var(--*) de tokens.css
283. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6ee7b7" (linha 3044); use var(--*) de tokens.css
284. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#047857" (linha 3045); use var(--*) de tokens.css
285. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#451a03" (linha 3049); use var(--*) de tokens.css
286. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fcd34d" (linha 3050); use var(--*) de tokens.css
287. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b45309" (linha 3051); use var(--*) de tokens.css
288. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#500724" (linha 3055); use var(--*) de tokens.css
289. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f9a8d4" (linha 3056); use var(--*) de tokens.css
290. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#be185d" (linha 3057); use var(--*) de tokens.css
291. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#083344" (linha 3061); use var(--*) de tokens.css
292. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#67e8f9" (linha 3062); use var(--*) de tokens.css
293. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0e7490" (linha 3063); use var(--*) de tokens.css
294. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1e1b4b" (linha 3067); use var(--*) de tokens.css
295. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#a5b4fc" (linha 3068); use var(--*) de tokens.css
296. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#4338ca" (linha 3069); use var(--*) de tokens.css
297. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1e1b4b" (linha 3073); use var(--*) de tokens.css
298. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#c7d2fe" (linha 3074); use var(--*) de tokens.css
299. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#4338ca" (linha 3075); use var(--*) de tokens.css
300. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#172554" (linha 3080); use var(--*) de tokens.css
301. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#93c5fd" (linha 3081); use var(--*) de tokens.css
302. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1d4ed8" (linha 3082); use var(--*) de tokens.css
303. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#2e1065" (linha 3085); use var(--*) de tokens.css
304. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#d8b4fe" (linha 3086); use var(--*) de tokens.css
305. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6d28d9" (linha 3087); use var(--*) de tokens.css
306. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#064e3b" (linha 3090); use var(--*) de tokens.css
307. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#6ee7b7" (linha 3091); use var(--*) de tokens.css
308. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#047857" (linha 3092); use var(--*) de tokens.css
309. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#451a03" (linha 3095); use var(--*) de tokens.css
310. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#fcd34d" (linha 3096); use var(--*) de tokens.css
311. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#b45309" (linha 3097); use var(--*) de tokens.css
312. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#500724" (linha 3100); use var(--*) de tokens.css
313. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#f9a8d4" (linha 3101); use var(--*) de tokens.css
314. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#be185d" (linha 3102); use var(--*) de tokens.css
315. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#083344" (linha 3105); use var(--*) de tokens.css
316. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#67e8f9" (linha 3106); use var(--*) de tokens.css
317. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#0e7490" (linha 3107); use var(--*) de tokens.css
318. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1e1b4b" (linha 3110); use var(--*) de tokens.css
319. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#a5b4fc" (linha 3111); use var(--*) de tokens.css
320. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#4338ca" (linha 3112); use var(--*) de tokens.css
321. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#1e1b4b" (linha 3115); use var(--*) de tokens.css
322. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#c7d2fe" (linha 3116); use var(--*) de tokens.css
323. **[aidd-linter-design]** (serious) `internal/studio/static/style.css`: Cor literal hardcoded "#4338ca" (linha 3117); use var(--*) de tokens.css
324. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#8ab4f8" (linha 57); use var(--*) de tokens.css
325. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#8ab4f8" (linha 57); use var(--*) de tokens.css
326. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#8ab4f8" (linha 57); use var(--*) de tokens.css
327. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#f28b82" (linha 58); use var(--*) de tokens.css
328. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#f28b82" (linha 58); use var(--*) de tokens.css
329. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#f28b82" (linha 58); use var(--*) de tokens.css
330. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#fdd663" (linha 59); use var(--*) de tokens.css
331. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#fdd663" (linha 59); use var(--*) de tokens.css
332. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#fdd663" (linha 59); use var(--*) de tokens.css
333. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#81c995" (linha 60); use var(--*) de tokens.css
334. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#81c995" (linha 60); use var(--*) de tokens.css
335. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#81c995" (linha 60); use var(--*) de tokens.css
336. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#c58af9" (linha 61); use var(--*) de tokens.css
337. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#c58af9" (linha 61); use var(--*) de tokens.css
338. **[aidd-linter-design]** (serious) `internal/studio/static/swarm.css`: Cor literal hardcoded "#c58af9" (linha 61); use var(--*) de tokens.css

## Próximo Passo Obrigatório (Governança até BUILD)

Este relatório é o insumo oficial para a **entrevista socrática de afunilamento** das correções.
Após a entrevista, gere e encadeie em ordem estrita os artefatos de correção:
1. `INTERVIEW` (decisões de correção por achado)
2. `INTENT` -> `SPEC` -> `PLAN` -> `ROUTE` -> `DAG` -> `SWARM`
3. Ciclo `BUILD` baseado neste relatório com reexecução do `legoloop-web-qa` como gate final.
