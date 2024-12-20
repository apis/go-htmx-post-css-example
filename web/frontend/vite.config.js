import cssnano from "cssnano"
import postcssImport from "postcss-import"
import postcssInherit from "postcss-inherit"


export default {
    base: process.env.NODE_ENV === "production" ? "/ui" : "/ui",
    css: {
        postcss: {
            plugins: [
                cssnano({preset: "default"}),
                postcssImport(),
                postcssInherit(),
            ],
        }
    },
}